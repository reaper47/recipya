---
title: Docker
weight: 3
---

A Docker image named [reaper99/recipya](https://hub.docker.com/layers/reaper99/recipya/nightly/images/sha256-b2238a11a53982953df5bbcfd7796a19fa382abf75d316b62fa05ac1c867332c?context=repo)
is produced nightly.

You can either install the application using [Docker](https://www.docker.com/) or
[Docker Compose](https://docs.docker.com/compose/).

## Using Docker

You first have to fetch it.

```bash
docker pull reaper99/recipya:nightly
```

Then, run the image. You must pass your [config.json](/guide/docs/installation/config-file) file to the container.

```bash
docker run -v path/to/config.json:/app/config.json -p [host port]:[port specified in config.json] -d reaper99/recipya:nightly reaper99/recipya:nightly
```

## Using Docker Compose

You can use Docker Compose to run the container. First, you need to modify the ports and the path to your local
[config.json](/guide/docs/installation/config-file) in the [compose.yaml](https://github.com/reaper47/recipya/blob/main/deploy/compose.yaml).
Then, start the application.

```bash
docker-compose up -d
```

Access the app through your browser at `http://localhost:[host's port]`.

If you are using Windows and you intend to access the app on other devices within your home network, please ensure to `Allow the connection` of the `Docker Desktop Backend`
inbound Windows Defender Firewall rule.