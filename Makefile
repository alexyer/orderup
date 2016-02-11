SOURCES=$(shell find . -name '*.go')

.PHONY: build
build:
	go get github.com/boltdb/bolt
	go get github.com/gorilla/mux
	go build -o orderup ${SOURCES}

.PHONY: install
install:
	go install
