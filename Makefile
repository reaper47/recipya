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

release:
ifdef tag
		go run ./releases/main.go -package github.com/reaper47/recipya -tag $(tag)
		gh release create $(tag) ./release/$(tag)/*
else
		@echo 'Add the tag argument, i.e. `make release tag=v1.0.0`'
endif

test:
	go test ./...

%:
	@:

.PHONY: release build run test