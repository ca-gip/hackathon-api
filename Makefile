.PHONY: build

REPO= github.com/ca-gip/hackathon-api
IMAGE= hackathon-api
TAG= dev
DOCKER_REPO= cagip

dependency:
	go mod vendor

build:
	GO111MODULE="on" CGO_ENABLED=0 go build -ldflags="-s" -v -o ./build/ $(GOPATH)/src/$(REPO)/main.go

darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o hackathon-api  $(GOPATH)/src/$(REPO)/main.go

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o  $(GOPATH)/src/$(REPO)/main.go

image:
	docker build -t "$(DOCKER_REPO)/$(IMAGE):$(TAG)" . --build-arg https_proxy=http://192.168.2.1:3128 --build-arg http_proxy=http://192.168.2.1:3128
	docker push "$(DOCKER_REPO)/$(IMAGE):$(TAG)" --build-arg https_proxy=http://192.168.2.1:3128 --build-arg http_proxy=http://192.168.2.1:3128

release:
	docker build -t "$(DOCKER_REPO)/$(IMAGE):$(TAG)" .
	docker push "$(DOCKER_REPO)/$(IMAGE):$(TAG)"

dep:
	glide install

