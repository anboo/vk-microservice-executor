package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vorkytaka/easyvk-go/easyvk"
	"io/ioutil"
	"log"
	"net/http"
	"flag"
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

type Executor struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Weight int `json:"weight"`
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

	w.Write([]byte("OK"))
}

func main() {
	masterHost := flag.String("master-host", "", "Master host for registration current executor")
	localHost := flag.String("host", "", "Current public host")
	localPort := flag.Int("port", 8000, "Choose public port")
	masterAuthorizationCode := flag.String("authorization", "", "Authorization header value")
	vkLogin := flag.String("login", "", "User login")
	vkPassword := flag.String("password", "", "User password")
	vkClient := flag.String("client-id", "", "VK OAuth client id")
	vkScope := flag.String("scope", "", "OAuth2 scopes")

	flag.Parse()

	log.Println("Try registration on master host " + *masterHost)
	regData, _ := json.Marshal(Executor{
		"http://" + *localHost,
		*localPort,
		1,
	})

	req, _ := http.NewRequest("POST", *masterHost + "/register-executor", bytes.NewBuffer(regData))
	req.Header.Add("Authorization", *masterAuthorizationCode)
	client := &http.Client{}
	res, err := client.Do(req); if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	log.Println(string(body))
	log.Println("Successfully registered in master host")

	vk, _ = easyvk.WithAuth(*vkLogin, *vkPassword, *vkClient, *vkScope)

	router := mux.NewRouter()
	router.HandleFunc("/request", DoRequest).Methods("POST")
	router.HandleFunc("/_health", HealthCheck).Methods("GET")

	log.Fatal(http.ListenAndServe(":8001", router))
}
