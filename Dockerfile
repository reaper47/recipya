# STAGE 1: prepare
FROM golang AS prepare

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /source 

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN apt-get update && apt-get install -y nodejs npm python3-venv python3-pip tesseract-ocr
RUN npm install -g yarn && yarn global add @vue/cli

# STAGE 2: build
FROM prepare as build

COPY . .

RUN go build -o main .
RUN cd /source/web && yarn install && yarn build

WORKDIR /dist 

RUN cp /source/main .  

# STAGE 3: run
FROM python:3.8-slim-buster as run

COPY --from=build /source/tools/requirements.txt /tools/requirements.txt
RUN pip install -r /tools/requirements.txt

COPY --from=build /dist/main /
COPY --from=build /source/tools/scraper.py /tools/scraper.py
COPY --from=build /source/web/dist /dist
COPY --from=build /source/dist/config.yaml /

EXPOSE 3001

ENTRYPOINT ["/main", "serve"]
