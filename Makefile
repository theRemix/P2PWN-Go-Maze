BINARY=p2p-ready
RUN_ENV=$(shell cat .env | xargs)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGEN=$(GOCMD) generate

all: build

build:
				$(GOGEN)
				$(GOBUILD) -o $(BINARY) -v

clean:
				$(GOCLEAN)
				rm -f $(BINARY)

run:
				$(GOBUILD) -o $(BINARY) -v
				env $(RUN_ENV) ./$(BINARY)

deps:
				$(GOGET) -d ./...

env:
				cp .env.sample .env
