NAME     := make-hashids
VERSION  := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: dep
dep:
	dep ensure

run:
	go run *.go

build: $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: test
test:
	go test

