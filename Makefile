default: all

all: clean build test

build:
	go build cmd/clouddriver/clouddriver.go

clean:
	go clean
	-rm ./clouddriver

run: clean build test
	./clouddriver

test:
	ginkgo -r

tools:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

vendor:
	go mod vendor

.PHONEY: all clean build run test tools
