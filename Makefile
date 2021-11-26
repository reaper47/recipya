.ONESHELL:

TAILWIND=pnpm dlx tailwindcss --jit -m --purge="./views/**/*.gohtml"

build-go:
	go build -ldflags="-s -w" -o dist/recipya main.go

build-css:
	 ${TAILWIND} -o static/css/tailwind.css 

watch-css: 
	${TAILWIND} -o static/css/tailwind.css -w

watch-js:
	cd views/tailwind && pnpm run watch


run:
	go run main.go serve

build: build-css build-go 
