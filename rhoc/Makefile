MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
LOCAL_BIN_PATH := ${PROJECT_PATH}/bin

ifeq ($(CLI_VERSION),)
CLI_VERSION := $(shell git rev-parse --short HEAD)
endif

GO_LDFLAGS := -X github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/internal/build.Version=$(CLI_VERSION) $(GO_LDFLAGS)


SHELL = bash

GO := go
GOFMT := gofmt

SHELL = bash

.PHONY: all
all: build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build: fmt vet
	go build $(BUILDFLAGS) -ldflags "${GO_LDFLAGS}" -o bin/rhoc cmd/rhoc/main.go

.PHONY: run
run: fmt vet 
	go run cmd/rhoc/main.go

.PHONY: deps
deps:
	go mod tidy

openapi/generate: openapi/generate/admin
.PHONY: openapi/generate

openapi/generate/admin: openapi-generator
	$(OPENAPI_GENERATOR) validate -i etc/openapi/connector_mgmt-private-admin.yaml
	$(OPENAPI_GENERATOR) generate -i etc/openapi/connector_mgmt-private-admin.yaml \
		-g go \
		-o pkg/api/admin \
		--package-name admin  \
		--additional-properties=enumClassPrefix=true \
		--additional-properties=isGoSubmodule=false \
		--additional-properties=withGoCodegenComment=true \
		-t etc/openapi/templates \
		--ignore-file-override etc/openapi/.openapi-generator-ignore

	$(GOFMT) -w pkg/api/admin
	
	rm pkg/api/admin/go.sum
	rm pkg/api/admin/go.mod
	rm pkg/api/admin/git_push.sh
	rm pkg/api/admin/.travis.yml
	rm pkg/api/admin/.openapi-generator-ignore
	rm pkg/api/admin/README.md

	rm -rf pkg/api/admin/docs
	rm -rf pkg/api/admin/api
	rm -rf pkg/api/admin/.openapi-generator

.PHONY: openapi/generate/admin

openapi/generate/public: openapi-generator
	$(OPENAPI_GENERATOR) validate -i etc/openapi/connector_mgmt.yaml
	$(OPENAPI_GENERATOR) generate -i etc/openapi/connector_mgmt.yaml \
		-g go \
		-o pkg/api/public \
		--package-name public  \
		--additional-properties=enumClassPrefix=true \
		--additional-properties=isGoSubmodule=false \
		--additional-properties=withGoCodegenComment=true \
		-t etc/openapi/templates \
		--ignore-file-override etc/openapi/.openapi-generator-ignore

	$(GOFMT) -w pkg/api/admin
	
	rm pkg/api/public/go.sum
	rm pkg/api/public/go.mod
	rm pkg/api/public/git_push.sh
	rm pkg/api/public/.travis.yml
	rm pkg/api/public/.openapi-generator-ignore
	rm pkg/api/public/README.md

	rm -rf pkg/api/public/docs
	rm -rf pkg/api/public/api
	rm -rf pkg/api/public/.openapi-generator

.PHONY: openapi/generate/public

OPENAPI_GENERATOR ?= ${LOCAL_BIN_PATH}/openapi-generator
NPM ?= "$(shell which npm)"
openapi-generator:
ifeq (, $(shell which ${NPM} 2> /dev/null))
	@echo "npm is not available please install it to be able to install openapi-generator"
	exit 1
endif
ifeq (, $(shell which ${LOCAL_BIN_PATH}/openapi-generator 2> /dev/null))
	@{ \
	set -e ;\
	mkdir -p ${LOCAL_BIN_PATH} ;\
	mkdir -p ${LOCAL_BIN_PATH}/openapi-generator-installation ;\
	cd ${LOCAL_BIN_PATH} ;\
	${NPM} install --prefix ${LOCAL_BIN_PATH}/openapi-generator-installation @openapitools/openapi-generator-cli ;\
	ln -s openapi-generator-installation/node_modules/.bin/openapi-generator-cli openapi-generator ;\
	}
endif


.PHONY: openapi-generator
