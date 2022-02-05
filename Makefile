.PHONY: build

REPO= github.com/ca-gip/hackathon-api
IMAGE= hackathon-api
TAG= dev
DOCKER_REPO= cagip

dependency:
	go mod vendor

build: dependency
	GO111MODULE="on" CGO_ENABLED=0 go build -ldflags="-s" -v -o ./build/hackathon-api $(GOPATH)/src/$(REPO)/main.go

image:
	docker buildx build --no-cache --push --platform linux/amd64,linux/arm64 -t "$(DOCKER_REPO)/$(IMAGE):$(TAG)" -f Dockerfile .

dep:
	glide install

