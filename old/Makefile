docker:
	go mod tidy -v
	go mod vendor
	docker build . -t reaper99/recipya:$(tag)

.PHONY: docker