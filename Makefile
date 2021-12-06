.ONESHELL:

TAILWIND=pnpm dlx tailwindcss --jit -m --purge="./views/**/*.gohtml"

build-go:
	go build -ldflags="-s -w" -o dist/recipya main.go

build-go-deploy:
	/usr/local/go/bin/go build -ldflags="-s -w" -o dist/recipya main.go

build-js:
	cd views/tailwind && pnpm i && pnpm run build

build-css:
	 ${TAILWIND} -o static/css/tailwind.css 

db-reset:
	go run main.go migrate down && go run main.go migrate up

run:
	go run main.go serve

watch-css: 
	${TAILWIND} -o static/css/tailwind.css -w

watch-js:
	cd views/tailwind && pnpm run watch

build: build-css build-js build-go

build-deploy: build-css build-js build-go-deploy
