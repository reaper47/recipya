# Build the application from source
FROM golang:1.23.0-alpine AS build-stage

WORKDIR /app

RUN apk add --update nodejs npm gcc g++ musl-dev git
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
ENV PATH="/go/bin:${PATH}"

COPY ../go.mod ../go.sum ./
RUN go mod download && go mod verify

COPY ./.git ./.git
COPY ./internal ./internal
COPY ./main.go ./
COPY ./web ./web

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/recipya main.go

# Deploy the application binary into a lean image
FROM alpine AS build-release-stage
LABEL org.opencontainers.image.source="https://github.com/reaper47/recipya"

WORKDIR /

RUN apk add --no-cache ffmpeg
RUN ffmpeg -version

COPY --from=build-stage /app/recipya /app/recipya

EXPOSE 8080

RUN chmod -R g+w /app && \
    chown -R nobody /app
USER nobody

ENTRYPOINT ["/app/recipya", "s"]