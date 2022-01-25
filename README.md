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
 export MONGOURI="mongodb://127.0.0.1:27017/myFirstDatabese?retryWrites=true&w=majority"
```

Install dependencies for pdf generation

```bash
 apt-get install xvfb libfontconfig wkhtmltopdf
```

Launch application

```bash
 go mod tidy
 go run main.go
```

TODO Dockerfile