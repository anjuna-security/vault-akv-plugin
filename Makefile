GOARCH = amd64

UNAME = $(shell uname -s)

PLUGIN_NAME=azure-key-vault

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/$(PLUGIN_NAME) cmd/vault_akv_plugin/main.go

test:
	GOOS=$(OS) GOARCH="$(GOARCH)" go test -v

coverage.out:
	GOOS=$(OS) GOARCH="$(GOARCH)" go test -coverprofile=coverage.out

cover: coverage.out
	GOOS=$(OS) GOARCH="$(GOARCH)" go tool cover -html=$<

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

enable:
	vault secrets enable -path=anjuna/ $(PLUGIN_NAME) 

clean:
	rm -f ./vault/plugins/$(PLUGIN_NAME) ./testapp/testapp coverage.out

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable
