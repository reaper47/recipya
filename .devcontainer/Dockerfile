FROM mcr.microsoft.com/devcontainers/go:latest

RUN apt-get update && apt-get install -y curl git wget

RUN curl -fsSL https://deb.nodesource.com/setup_current.x | bash - && \
    apt-get install -y build-essential nodejs

RUN CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest