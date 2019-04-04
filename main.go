package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vorkytaka/easyvk-go/easyvk"
	"log"
	"net/http"
	"os"
)

type Response struct {
	RequestId string `json:"request_id"`
	Result interface{} `json:"result"`
	ResultAsString string `json:"result_as_string"`
}

type Request struct {
	Id string `json:"id"`
	Method string `json:"method"`
	Parameters map[string]string `json:"parameters"`
	Response Response `json:"response"`
}

var vk easyvk.VK

func DoRequest(w http.ResponseWriter, r *http.Request)  {
	var req Request
	_ = json.NewDecoder(r.Body).Decode(&req)

	resBytes, err := vk.Request(req.Method, req.Parameters); if err != nil {
		log.Println("Error" + err.Error())

		errRes := make(map[string]string)
		errRes["error"] = err.Error()

		json.NewEncoder(w).Encode(errRes)
		return
	}

	var response interface{}
	unmarshalErr := json.Unmarshal(resBytes, &response); if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}

	errSerialize := json.NewEncoder(w).Encode(response); if errSerialize != nil {
		log.Fatal(errSerialize)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, errSerialize := w.Write([]byte("OK")); if errSerialize != nil {
		log.Fatal(errSerialize)
	}
}

func main() {
	vk, _ = easyvk.WithAuth(os.Getenv("VK_EMAIL"), os.Getenv("VK_PASSWORD"), os.Getenv("VK_CLIENT"), os.Getenv("VK_SCOPE"))
	vk.AccessToken = ""

	log.Println("Try authorization in master host")

	router := mux.NewRouter()
	router.HandleFunc("/request", DoRequest).Methods("POST")
	router.HandleFunc("/_health", HealthCheck).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
