SOURCES=$(shell find . -name '*.go')

.PHONY: build
build:
	go get github.com/tools/godep
	godep restore
	go build -o orderup ${SOURCES}

.PHONY: install
install:
	go install ./orderup
