.ONESHELL:

TAILWIND=pnpm dlx tailwindcss --jit -m --purge="./views/**/*.gohtml"

build-api:
	go build -o dist/recipya main.go

build-css:
	 ${TAILWIND} -o static/css/tailwind.css 

build-go:
	go build -ldflags="-s -w"

watch: 
	${TAILWIND} -o static/css/tailwind.css -w

run:
	go run main.go serve

build: build-css build-api 
