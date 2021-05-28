FROM golang:1.16.4-alpine3.13 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build 

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist 

RUN cp /build/main .  

FROM scratch

COPY --from=builder /dist/main /
COPY ./build/config.yaml /

EXPOSE 3001

ENTRYPOINT ["/main", "serve"]
