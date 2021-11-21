.ONESHELL:

TAILWIND=pnpm dlx tailwindcss --jit -m --purge="./views/**/*.gohtml"

build-api:
	go build -o dist/recipya main.go
	mkdir -p dist/views
	cp -R views/layouts views/*.gohtml dist/views/

build-css:
	 ${TAILWIND} -o static/css/tailwind.css 

build-go:
	go build -ldflags="-s -w"

watch: 
	${TAILWIND} -o static/css/tailwind.css -w
	
run:
	go run main.go serve

build: build-css build-api 
