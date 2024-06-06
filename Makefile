default: all

all: clean lint build test

build:
	go build cmd/clouddriver/clouddriver.go

clean:
	go clean
	-rm ./clouddriver

lint:
	golangci-lint run

run: clean lint build test
	./clouddriver

test:
	ginkgo -r

tools:
	go get github.com/onsi/ginkgo/v2/ginkgo
	go get github.com/onsi/gomega/...

vendor:
	go mod vendor

.PHONY: all clean build lint run test tools
