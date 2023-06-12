.ONESHELL:

build:
	go generate ./...
	go build -ldflags="-s -w" -o bin/recipya main.go

release:
ifdef tag
		sh build.sh github.com/reaper47/recipya $(tag)
		gh release create $(tag) ./release/$(tag)/*
else
		@echo 'Add the tag argument, i.e. `make release tag=v1.0.0`'
endif

test:
	go test ./...

%:
	@:

.PHONY: release build run test