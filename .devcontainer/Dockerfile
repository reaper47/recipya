FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

RUN apt-get update && apt-get install -y curl wget git

RUN curl -fsSL https://deb.nodesource.com/setup_current.x | bash - && \
    apt-get install -y nodejs \
    build-essential

RUN CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/a-h/templ/cmd/templ@latest