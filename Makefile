.ONESHELL:

ifeq ($(OS),Windows_NT)
    TARGET_EXT = .exe
else
    TARGET_EXT =
endif

TARGET = recipya$(TARGET_EXT)

build:
	go generate ./...
	go build -ldflags="-s -w" -o bin/$(TARGET) main.go

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

release:
ifdef tag
		go run ./releases/main.go -package github.com/reaper47/recipya -tag $(tag)
		gh release create $(tag) ./releases/$(tag)/*
else
		@echo 'Add the tag argument, i.e. `make release tag=v1.0.0`'
endif

sponsors:
	cd ./scripts/sponsors
	npm install sponsorkit@0.6.1
	npx sponsorkit -o ../../docs/website/static/img/

test:
	go test ./...

%:
	@:

.PHONY: release build run sponsors test