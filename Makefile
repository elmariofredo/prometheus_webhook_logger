APP_NAME  := prometheus_webhook_logger
DEBUG_POSTFIX := -debugger
NAMESPACE := sysincz
IMAGE := $(NAMESPACE)/$(APP_NAME)
REGISTRY := docker.io
ARCH := linux 
VERSION := $(shell git describe --tags 2>/dev/null)
ifeq "$(VERSION)" ""
VERSION := $(shell git rev-parse --short HEAD)
endif
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE=$(shell date +%FT%T%z)
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Branch=$(BRANCH) -X main.BuildDate=$(BUILD_DATE)"

.PHONY: clean

clean:
	rm -rf bin/%/$(APP_NAME)

dep:
	go get -v ./...

build: clean dep
	GOOS=$(ARCH) GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -a -installsuffix cgo -o bin/$(APP_NAME) ./

run: build 
	./bin/$(APP_NAME) --config ./logger.yaml

send-alert:
	curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" --data "@alert_example.json" localhost:9099/alerts

docker: build build-logger-image

build-logger-image: build
	docker build . -t $(IMAGE):$(VERSION)


docker-push: docker
	docker push $(IMAGE):$(VERSION)


webhook: build-logger-image
	docker run -p 9099:9099/tcp --network host -v "/tmp/config:/config" --rm --name $(APP_NAME) $(IMAGE):$(VERSION)

debug_webhook: build-logger-image
	docker run -it -p 9099:9099/tcp --network host --rm --entrypoint "/bin/bash" --name $(APP_NAME) $(IMAGE):$(VERSION) 
