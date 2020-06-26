GO=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
TAG=latest
BIN=labelgun
BIN_PATH=bin/$(BIN)
IMAGE=dailyhotel/$(BIN)

all: image
	docker push $(IMAGE):$(TAG)

build:
	$(GO) build -a -installsuffix cgo -o $(BIN_PATH) .

image: build
	docker build \
		--build-arg VCS_REF=`git rev-parse --short HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t $(IMAGE):$(TAG) .

.PHONY: clean

clean:
	rm $(BIN_PATH)
