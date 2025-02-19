# Build the application from source
FROM golang:1.24.0-alpine AS build-stage

WORKDIR /app

RUN apk add --update g++ gcc git musl-dev nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
ENV PATH="/go/bin:${PATH}"

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY ./.git ./.git
COPY ./internal ./internal
COPY ./main.go ./
COPY ./web ./web

RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/recipya main.go

# Deploy the application binary into a lean image
FROM alpine AS build-release-stage
LABEL org.opencontainers.image.source="https://github.com/reaper47/recipya"
LABEL org.opencontainers.image.description="A clean, simple and powerful recipe manager your whole family will enjoy."
LABEL org.opencontainers.image.licenses="GPL-3.0"
WORKDIR /app

RUN apk add --no-cache ffmpeg
RUN ffmpeg -version

RUN adduser -D recipya
USER recipya:recipya
RUN mkdir -p /home/recipya/.config/Recipya
COPY --from=build-stage /app/recipya .

EXPOSE 8078

ENTRYPOINT ["/app/recipya", "serve"]