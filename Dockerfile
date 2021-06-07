# STAGE 1: prepare
FROM golang AS prepare

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /source 

COPY vendor .

RUN apt-get update 
RUN apt-get install -y curl git wget unzip libgconf-2-4 gdb libstdc++6 libglu1-mesa fonts-droid-fallback lib32stdc++6 python3
RUN apt-get clean

RUN git clone https://github.com/flutter/flutter.git /usr/local/flutter

ENV PATH="/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin:${PATH}"

RUN flutter doctor -v
RUN flutter channel master
RUN flutter upgrade
RUN flutter config --enable-web

# STAGE 2: build
FROM prepare as build

COPY . .

RUN go build -o main .
RUN cd /source/web/app && flutter build web

WORKDIR /dist 

RUN cp /source/main .  

# STAGE 3: run
FROM scratch as run

COPY --from=build /dist/main /
COPY --from=build /source/web/app/build/web /web
COPY --from=build /source/build/config.yaml /

EXPOSE 3001

ENTRYPOINT ["/main", "serve"]
