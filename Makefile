GO 			 ?= /usr/local/go/bin/go
DOCKER 	 ?= docker
DEP			 ?= $(GOPATH)/bin/dep

IMAGE_REPO ?= bradhe
IMAGE_NAME  = blobd
IMAGE_TAG  ?= latest
IMAGEPATH		= $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

BINDIR 	= .
BINNAME = blobd
BINPATH	= $(BINDIR)/$(BINNAME)

setup:
	@if [ ! -f $(DEP) ]; then echo "Install godep https://github.com/golang/dep"; exit 1; fi

clean:
	rm -rf ./server/ui/build
	rm -rf ./server/ui/.blobd.last-build.checksum
	rm -rf ./blobd

generate:
	$(GO) generate ./...

build: generate
	$(GO) build -o $(BINPATH) ./cmd/blobd

build_linux: generate
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINPATH) ./cmd/blobd

test: generate
	$(GO) test -tags 'integration unit' ./...	

images: build_linux
	$(DOCKER) build -t $(IMAGEPATH) .

release: images
	$(DOCKER) push $(IMAGEPATH)

serve: build
	docker-compose up --build
