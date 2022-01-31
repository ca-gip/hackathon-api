# Hackathon-Api

### A REST Api to donate all your extra Cryptocurrency for to the banking cause

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

### Call API example

To create a new donation

```bash
curl --request POST \
    --url http://localhost:8080/donation \
    --header 'Content-Type: application/json' \
    --data '{
    "donatorName": "Rob STRACK",
    "amount": 125,
    "moneyType": "€"
    }'
```

To get a donation

```bash
curl --request GET \
    --url 'http://localhost:8080/donation/<DONATION_ID>' \
    --header 'Content-Type: application/json'
```

To get all donation

```bash
curl --request GET \
    --url 'http://localhost:8080/donations' \
    --header 'Content-Type: application/json'
```



# Run on ARM
docker run  --platform linux/arm64 -e "MONGOURI=mongodb://192.168.1.68:27017/donation?retryWrites=true&w=majority" -d --name hackathon cagip/hackathon-api:dev