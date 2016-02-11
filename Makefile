SOURCES=$(shell find ./orderup -name '*.go')

.PHONY: build
build:
	go get github.com/boltdb/bolt
	go get github.com/gorilla/mux
	go build -o orderup-server ${SOURCES}

.PHONY: install
install:
	go install ./orderup
