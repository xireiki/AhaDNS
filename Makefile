NAME = ahadns
COMMIT = $(shell git rev-parse --short HEAD)

GOHOSTOS = $(shell go env GOHOSTOS)
GOHOSTARCH = $(shell go env GOHOSTARCH)

PARAMS = -v -trimpath
MAIN = .
PREFIX ?= $(shell go env GOPATH)

.PHONY : build clean build_all

OSList = linux windows darwin android freebsd
ArchList = arm64 amd64

build:
	go build $(PARAMS) $(MAIN)

install:
	go build -o $(PREFIX)/bin/$(NAME) $(PARAMS) $(MAIN)

build_all:
	@for os in $(OSList); do \
		for arch in $(ArchList); do \
			echo "Building for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build $(PARAMS) -o $(NAME)-$$os-$$arch $(MAIN); \
		done; \
	done

clean:
	@rm -rf ahadns-*
