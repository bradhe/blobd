GO 			 ?= /usr/local/go/bin/go
GOVENDOR ?= $(GOPATH)/bin/govendor
DOCKER 	 ?= docker

IMAGE_REPO ?= bradhe
IMAGE_NAME  = blobd
IMAGE_TAG  ?= latest

setup:
	go get -v github.com/kardianos/govendor

clean:
	rm ./cmd/blobd/blobd

build:
	$(GO) generate ./...
	$(GO) build -o ./cmd/blobd/blobd ./cmd/blobd

build_linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o ./cmd/blobd/blobd ./cmd/blobd

test: build
	$(GOVENDOR) test -tags 'integration unit' ./...	

images: build_linux
	$(DOCKER) build -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) ./cmd/blobd

release: images
	$(DOCKER) push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

