# Hackathon-Api

### A REST Api to donate all your extra Cryptocurrency for the banking cause

This repository shows the source code for building application with Golang using the Gin-gonic framework and MongoDB.

## Prerequisites

First you need a MongoDB database, here an example with Docker

```bash
 docker run -p 27017:27017 --name some-mongo -d mongo
```

You need to set Mongo URI env var 

```bash
 export MONGOURI="mongodb://127.0.0.1:27017/donation?retryWrites=true&w=majority"
```

Install dependencies for pdf generation

```bash
 apt-get install xvfb libfontconfig wkhtmltopdf
```

Launching application

```bash
 go mod tidy
 go run main.go
```

### To build the docker image

You need to change the repositories address in Makefile ``` DOCKER_REPO= your_repo ```

If you need to be authenticated use ``` docker login ``` command with your credentials

```bash
 make image
```

Run image

```bash
# On arm64
docker run  --platform linux/arm64 -e "MONGOURI=mongodb://<YOUR_IP>:27017/donation?retryWrites=true&w=majority" -p 8080:8080 -d --name hackathon cagip/hackathon-api:dev

# On amd64
docker run  --platform linux/amd64 -e "MONGOURI=mongodb://<YOUR_IP>:27017/donation?retryWrites=true&w=majority" -p 8080:8080 -d --name hackathon cagip/hackathon-api:dev

```

### API Call example

To create a new donation

```bash
curl --request POST \
    --url http://localhost:8080/donation \
    --header 'Content-Type: application/json' \
    --data '{
    "donatorName": "Rob STRACK",
    "amount": 125,
    "moneyType": "BTC"
    }'
```

To get a donation

```bash
curl --request GET \
    --url 'http://localhost:8080/donation/<DONATION_ID>' \
    --header 'Content-Type: application/json'
```

To download a reward

```bash
curl --request GET \
    --url 'http://localhost:8080/document/<PDF_REF>' \
    --header 'Content-Type: application/json'
```


To get all donation

```bash
curl --request GET \
    --url 'http://localhost:8080/donations' \
    --header 'Content-Type: application/json'
```
