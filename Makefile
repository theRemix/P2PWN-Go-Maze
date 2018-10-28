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
				env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY)-osx64 -v
				cd $(BIN_PATH) zip $(BINARY)-osx64.zip $(BINARY)-osx64

build-win64:
				@echo go-gl and glfw are not yet supported on this platform
				@# env GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY)-win64 -v
				@# cd $(BIN_PATH) zip $(BINARY)-win64.zip $(BINARY)-win64

build-linux64:
				@echo go-gl and glfw are not yet supported on this platform
				@# env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(BINARY)-linux64 -v
				@# cd $(BIN_PATH) zip $(BINARY)-linux64.zip $(BINARY)-linux64

release: clean generate build-osx64 build-win64 build-linux64

clean:
				$(GOCLEAN)
				rm -f $(BIN_PATH)/*

run:
				$(GOBUILD) -o $(BIN_PATH)/$(BINARY)-osx64 -v
				env $(RUN_ENV) $(BIN_PATH)/$(BINARY)-osx64

generate:
				$(GOGEN)

deps:
				$(GOGET) -d ./...

env:
				cp .env.sample .env
