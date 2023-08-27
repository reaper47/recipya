# Build the application from source
FROM golang:1.21.0-alpine AS build-stage

WORKDIR /app

RUN apk add --update nodejs npm

COPY ../go.mod ../go.sum ./
RUN go mod download && go mod verify

COPY ./internal ./internal
COPY ./main.go ./
COPY ./static ./static
COPY ./web ./web

RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/recipya main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM istio/distroless AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/recipya /app/recipya

EXPOSE 8080

USER root:root

ENTRYPOINT ["/app/recipya", "s"]