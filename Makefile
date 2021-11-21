.ONESHELL:

TAILWIND=pnpm dlx tailwindcss --jit -m --purge="./views/**/*.gohtml"

build-go:
	go build -ldflags="-s -w" -o dist/recipya main.go

build-css:
	 ${TAILWIND} -o static/css/tailwind.css 

watch: 
	${TAILWIND} -o static/css/tailwind.css -w

run:
	go run main.go serve

build: build-css build-go 
