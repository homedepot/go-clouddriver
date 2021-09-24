default: all

all: clean lint build test

build:
	go build cmd/clouddriver/clouddriver.go

clean:
	go clean
	-rm ./clouddriver

lint:
	golangci-lint run \
		--enable misspell \
		--enable wsl \
		--print-issued-lines=false \
		--skip-files .*_test.go \
		--timeout=3m0s \
		--out-format=colored-line-number \
		--issues-exit-code=1 ./...

run: clean lint build test
	./clouddriver

test:
	ginkgo -r

tools:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...

vendor:
	go mod vendor

.PHONY: all clean build lint run test tools
