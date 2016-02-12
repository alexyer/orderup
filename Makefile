BOT_SOURCES=$(shell find ./cmd/orderup -name '*.go')
CLI_SOURCES=$(shell find ./cmd/ou -name '*.go')

.PHONY: build
build:
	go get github.com/boltdb/bolt
	go get github.com/gorilla/mux
	go get github.com/codegangsta/cli
	go build -o ./bin/orderup ${BOT_SOURCES}
	go build -o ./bin/ou ${CLI_SOURCES}

.PHONY: install
install:
	go install ./cmd/ou
	go install ./cmd/orderup
