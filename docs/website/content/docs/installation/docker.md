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

Then, run the image. The `-e` environment variables are described below.

```bash
docker run -d \
  --name recipya
  --restart unless-stopped
  -p <host port>:<port specified by RECIPYA_SERVER_PORT> \
  -v recipya-data:/root/.config/Recipya \
  -e RECIPYA_EMAIL=my@email.com \
  -e RECIPYA_EMAIL_SENDGRID=API_KEY \
  -e RECIPYA_VISION_KEY=KEY_1 \
  -e RECIPYA_VISION_ENDPOINT=https://{resource}.cognitiveservices.azure.com \
  -e RECIPYA_SERVER_AUTOLOGIN=false \
  -e RECIPYA_SERVER_IS_DEMO=false \
  -e RECIPYA_SERVER_IS_PROD=false \
  -e RECIPYA_SERVER_NO_SIGNUPS=false \
  -e RECIPYA_SERVER_PORT=8078 \
  -e RECIPYA_SERVER_URL=http://0.0.0.0 \
  reaper99/recipya:nightly
```

## Using Docker Compose

You can use Docker Compose to run the container. First, download the [compose.yaml](https://github.com/reaper47/recipya/blob/main/deploy/compose.yaml) file. 
Modify the `environment` and `ports` sections. The environment variables are described below. Then, start the application.

```bash
docker-compose up -d
```

Access the app through your browser at `http://localhost:[host port]`.

If you are using Windows and you intend to access the app on other devices within your home network, please ensure to `Allow the connection` of the `Docker Desktop Backend`
inbound Windows Defender Firewall rule.

## Environment Variables

| Variable                  | Description                                                                                                                                           |       Required       |
|---------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|:--------------------:|
| RECIPYA_EMAIL             | The administrator’s email address. It is usually the email address of your [SendGrid](https://sendgrid.com/) account.                                 |   {{< icon "x" >}}   |
| RECIPYA_EMAIL_SENDGRID    | Your [SendGrid](https://app.sendgrid.com/settings/api_keys) API key. The free tier should be sufficient for your needs.                               |   {{< icon "x" >}}   |
| RECIPYA_SERVER_AUTOLOGIN  | Whether to login automatically into the application. Set to `true` when you don't need user accounts.                                                 |   {{< icon "x" >}}   |
| RECIPYA_SERVER_IS_DEMO    | Whether the app is a demo version. Its value can be either `true` or `false`.                                                                         |   {{< icon "x" >}}   |
| RECIPYA_SERVER_IS_PROD    | Whether the app is in production. Its value can be either `true` or `false`.                                                                          |   {{< icon "x" >}}   |
| RECIPYA_SERVER_NO_SIGNUPS | Whether to disable user account registrations. Set to `true` when you don't want people to create accounts.                                           |   {{< icon "x" >}}   |
| RECIPYA_SERVER_PORT       | The port the app will be served through if localhost.                                                                                                 | {{< icon "check" >}} |
| RECIPYA_SERVER_URL        | The website the app is served on. This URL will serve as the base link in the emails. <br/>                                                           | {{< icon "check" >}} |
| RECIPYA_VISION_KEY        | The **KEY 1** variable displayed on the Keys and endpoint tab of your Computer vision resource in the [Azure Portal](https://portal.azure.com/#home). |   {{< icon "x" >}}   |
| RECIPYA_VISION_ENDPOINT   | The Endpoint variable displayed on the Keys and endpoint tab of your Computer vision resource in the [Azure Portal](https://portal.azure.com/#home).  |   {{< icon "x" >}}   |
