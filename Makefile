GO 			 ?= /usr/local/go/bin/go
GOVENDOR ?= $(GOPATH)/bin/govendor
DOCKER 	 ?= /usr/local/bin/docker

IMAGE_REPO ?= bradhe
IMAGE_NAME  = blobd
IMAGE_TAG  ?= latest

build:
	$(GO) build -a -o ./cmd/blobd/blobd ./cmd/blobd

test:
	$(GOVENDOR) test -tags 'integration unit' ./...	

images: build
	$(DOCKER) build -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) ./cmd/blobd

