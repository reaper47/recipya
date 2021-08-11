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

RUN apt-get update && apt-get install -y nodejs npm python3-venv
RUN python3 -m venv /source/tools/venv && /source/tools/venv/bin/pip install --upgrade pip
RUN npm install -g yarn && yarn global add @vue/cli

# STAGE 2: build
FROM prepare as build

COPY . .

RUN go build -o main .
RUN /source/tools/venv/bin/pip install -r /source/tools/requirements.txt
RUN cd /source/web && yarn install && yarn build

WORKDIR /dist 

RUN cp /source/main .  

# STAGE 3: run
FROM gcr.io/distroless/python3 as run

COPY --from=build /dist/main /
COPY --from=build /source/tools/scraper/scraper.py /tools/scraper.py
COPY --from=build /source/tools/venv /tools/venv
COPY --from=build /source/web/dist /dist
COPY --from=build /source/dist/config.yaml /

EXPOSE 3001

ENTRYPOINT ["/main", "serve"]
