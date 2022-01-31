FROM golang:latest
WORKDIR /root/
RUN apt-get update && \
      apt-get install -y --no-install-recommends \
      xvfb libfontconfig wkhtmltopdf \
      && rm -rf /var/lib/apt/lists/*
COPY . .
RUN go mod vendor
RUN GO111MODULE="on" CGO_ENABLED=0 go build -ldflags="-s" -v -o ./build/hackathon-api ./main.go

EXPOSE 8080
CMD ["./build/hackathon-api"]
