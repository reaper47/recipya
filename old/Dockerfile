# STAGE 1: prepare
FROM golang:1.17.0-bullseye AS prepare

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /source 

COPY go.mod .
COPY go.sum .
RUN go mod download

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y nodejs npm python3-venv python3-pip tesseract-ocr libtesseract-dev libleptonica-dev
RUN npm install -g yarn && yarn global add @vue/cli

# STAGE 2: build
FROM prepare as build

COPY . .

RUN go build -o main .
RUN cd /source/ui && yarn install && yarn build

WORKDIR /dist 

RUN cp /source/main .  

# STAGE 3: run
FROM python:3.8-slim-buster as run

RUN apt-get update && apt-get install -y tesseract-ocr libtesseract-dev libleptonica-dev

COPY --from=build /source/tools/requirements.txt /tools/requirements.txt
RUN pip install -r /tools/requirements.txt

COPY --from=build /dist/main /
COPY --from=build /source/data/blacklist_units.txt /data/blacklist_units.txt
COPY --from=build /source/data/fruits_veggies.txt /data/fruits_veggies.txt
COPY --from=build /source/tools/scraper.py /tools/scraper.py
COPY --from=build /source/ui/dist /ui/dist
COPY --from=build /source/dist/config.yaml /

EXPOSE 3001

ENTRYPOINT ["/main", "serve"]
