GO 			 ?= /usr/local/go/bin/go
GOVENDOR ?= $(GOPATH)/bin/govendor
DOCKER 	 ?= /usr/local/bin/docker

IMAGE_REPO ?= bradhe
IMAGE_NAME  = blobd
IMAGE_TAG  ?= latest

setup:
	go get -v github.com/kardianos/govendor

clean:
	rm ./cmd/blobd/blobd

build:
	$(GO) build -o ./cmd/blobd/blobd ./cmd/blobd

test:
	$(GOVENDOR) test -tags 'integration unit' ./...	

images: build
	$(DOCKER) build -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) ./cmd/blobd

release: images
	$(DOCKER) login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)
	$(DOCKER) push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

