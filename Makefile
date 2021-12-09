.ONESHELL:

build-go:
	go build -ldflags="-s -w" -o dist/recipya main.go

build-go-deploy:
	/usr/local/go/bin/go build -ldflags="-s -w" -o dist/recipya main.go

build-static:
	 cd views/tailwind && pnpm i && pnpm run build 

db-reset:
	go run main.go migrate down && go run main.go migrate up

run:
	go run main.go serve

watch-css: 
	cd views/tailwind && pnpm run watch-css 

watch-js:
	cd views/tailwind && pnpm run watch

build: build-static build-go

build-deploy: build-static build-go-deploy
