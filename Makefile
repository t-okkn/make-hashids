NAME     := make-hashids
VERSION  := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
DSTDIR  := /srv/http/osaka/bin
USER    := http
GROUP   := http
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

run:
	go run *.go

.PHONY: dep
dep:
	dep ensure

build: $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

install:
	\cp -r bin/$(NAME) $(DSTDIR)/
	chown $(USER):$(GROUP) $(DSTDIR)/$(NAME)

uninstall: revoke_service
	rm -f $(DSTDIR)/$(NAME)

create_service:
	echo -e "[Unit]\nDescription = $(NAME)(Golang App)\n\n[Service]\nEnvironment = \"GIN_MODE=release\"\nWorkingDirectory = $(DSTDIR)/\n\nExecStart = $(DSTDIR)/$(NAME)\nExecStop = /bin/kill -HUP $MAINPID\nExecReload = /bin/kill -HUP $MAINPID && $(DSTDIR)/$(NAME)\n\nRestart = always\nType = simple\n\n[Install]\nWantedBy = multi-user.target" | sudo tee /etc/systemd/system/$(NAME).service
	systemctl enable $(NAME)

revoke_service: /etc/systemd/system/$(NAME).service
	systemctl stop $(NAME)
	systemctl disable $(NAME)
	rm -f /etc/systemd/system/$(NAME).service

clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: test
test:
	go test

