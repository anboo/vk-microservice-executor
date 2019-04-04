# Installation
```bash
cp .env.dist .env
docker build -t anboo/golang-vk-proxy-executor .
docker run  -p 8888:8000 --env-file .env anboo/golang-vk-proxy-executor
```

# Usage

POST /request
```
{
    "id": "266283c3-caf0-47ac-baf0-a4a827edb77f",
    "method": "users.get",
    "parameters": {
        "user_ids": "31292206,31292206"
    }
}
```
