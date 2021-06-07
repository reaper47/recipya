docker:
	go mod tidy -v
	go mod vendor
	docker build . -t reaper99/recipe-hunter:$(tag)

.PHONY: docker