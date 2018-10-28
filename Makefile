BINARY=go-maze
BIN_PATH=bin
RUN_ENV=$(shell cat .env | xargs)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGEN=$(GOCMD) generate

all: build

build-osx64:
				env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY) -v

build-win64:
				@echo go-gl and glfw are not yet supported on this platform
				@# env GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY) -v

build-linux64:
				@echo go-gl and glfw are not yet supported on this platform
				@# env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY) -v

release: generate build-osx64 build-win64 build-linux64

clean:
				$(GOCLEAN)
				rm -f $(BIN_PATH)/*

run:
				$(GOBUILD) -o $(BIN_PATH)/$(BINARY) -v
				env $(RUN_ENV) $(BIN_PATH)/$(BINARY)

generate:
				$(GOGEN)

deps:
				$(GOGET) -d ./...

env:
				cp .env.sample .env
