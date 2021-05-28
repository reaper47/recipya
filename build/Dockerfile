FROM golang:1.16.4-alpine3.13

RUN apk add --no-cache git nano

WORKDIR /app/recipe-hunter

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

#RUN go build -o ./bin/recipe-hunter

EXPOSE 3001

CMD ["sh"]
