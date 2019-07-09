NAME     := make-hashids
VERSION  := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
DSTDIR  := /srv/http/osaka/bin
USER    := http
GROUP   := http
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: dep
dep:
	dep ensure

run:
	go run *.go

build: $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

install:
	\cp -r bin/$(NAME) $(DSTDIR)/
	chown $(USER):$(GROUP) $(DSTDIR)/$(NAME)

uninstall:
	rm -f $(DSTDIR)/$(NAME)

clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: test
test:
	go test

