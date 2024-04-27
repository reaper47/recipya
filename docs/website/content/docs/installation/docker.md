---
title: Docker
weight: 3
---

A Docker image named [reaper99/recipya](https://hub.docker.com/layers/reaper99/recipya/nightly/images/sha256-b2238a11a53982953df5bbcfd7796a19fa382abf75d316b62fa05ac1c867332c?context=repo)
is produced nightly whenever the `main` branch has new commits.

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
  -p 8085:8078 \
  -v recipya-data:/root/.config/Recipya \
  -e RECIPYA_SERVER_PORT=8078 \
  -e RECIPYA_SERVER_URL=http://0.0.0.0 \
  reaper99/recipya:nightly
```

Recipya can be accessed from your host machine at [http://localhost:8085](http://localhost:8085).

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

| Variable                  | Description                                                                                                                                                                                                                                  |
|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RECIPYA_DI_ENDPOINT       | The **Endpoint** variable displayed in the *Keys and endpoint* tab of your [Azure AI Document Intelligence](https://azure.microsoft.com/en-us/products/ai-services/ai-document-intelligence) resource in the Azure Portal.<br>Default: `""`. |
| RECIPYA_DI_KEY            | The **KEY 1** variable displayed in the *Keys and endpoint* tab of your *Document Intelligence* resource in the [Azure Portal](https://portal.azure.com/#home).<br>Default: `""`.                                                            |
| RECIPYA_EMAIL             | The email address of your [SendGrid](https://sendgrid.com/) account.<br>Default: `""`.                                                                                                                                                       |
| RECIPYA_EMAIL_SENDGRID    | Your [SendGrid](https://app.sendgrid.com/settings/api_keys) API key. The free tier should be sufficient for your needs.<br>Default: `""`.                                                                                                    |
| RECIPYA_SERVER_AUTOLOGIN  | Whether to login automatically into the application. Can be `true` or `false`.<br>Default: `false`.                                                                                                                                          |
| RECIPYA_SERVER_IS_DEMO    | Whether the app is a demo version. Can be `true` or `false`.<br>Default: `false`.                                                                                                                                                            |
| RECIPYA_SERVER_IS_PROD    | Whether the app is in production. Can be `true` or `false`.<br>Default: `false`.                                                                                                                                                             |
| RECIPYA_SERVER_NO_SIGNUPS | Whether to disable user account registrations. Set to `true` when you don't want people to create accounts.<br>Default: `false`.                                                                                                             |
| RECIPYA_SERVER_PORT       | The port the app will be served through if localhost.<br>**Is required**.                                                                                                                                                                    |
| RECIPYA_SERVER_URL        | The website the app is served on. This URL will serve as the base link in the emails.<br>**Is required**.                                                                                                                                    |

### Deprecations

The following table lists deprecated environment variables in v1.2.0. They will be removed in v1.3.0.

| Variable                  | Reason                                                                                                                                                                                                                                                         |
|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RECIPYA_VISION_ENDPOINT   | Replaced the use of [Azure AI Vision](https://azure.microsoft.com/en-us/products/ai-services/ai-vision) to digitize recipes in favor of [Azure AI Document Intelligence](https://azure.microsoft.com/en-us/products/ai-services/ai-document-intelligence).     |
| RECIPYA_VISION_KEY        | Same as above.                                                                                                                                                                                                                                                 |
