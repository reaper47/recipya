---
title: Docker
weight: 3
---

A [stable](https://hub.docker.com/layers/reaper99/recipya/v1.1.0/images/sha256-fb1457919f132ebf6969f9c155d81bb60b0d6b0b1610bc692259b6b9c287479e?context=repo) Docker
is produced on every release. An [unstable](https://hub.docker.com/layers/reaper99/recipya/nightly/images/sha256-b2238a11a53982953df5bbcfd7796a19fa382abf75d316b62fa05ac1c867332c?context=repo) one 
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

### Updating your container

Run the command below for a quick update with Watchtower. 

Remember to replace `recipya` with your actual container name if it differs and to [back up your volume data](#back-up-a-volume).

```bash
docker run --rm --volume /var/run/docker.sock:/var/run/docker.sock containrrr/watchtower --run-once recipya
```

## Using Docker Compose

You can use Docker Compose to run the container. First, download the [compose.yaml](https://github.com/reaper47/recipya/blob/main/deploy/compose.yaml) file. 
Modify the `environment` and `ports` sections. The environment variables are described below. Then, start the application.

```bash
docker-compose up -d
```

Access the app through your browser at `http://localhost:[host port]`.

If you are using Windows and you intend to access the app on other devices within your home network, please ensure to 
`Allow the connection` of the `Docker Desktop Backend` inbound Windows Defender Firewall rule.

### Updating your container

Follow these steps to update Recipya. Remember to [back up your volume data](#back-up-a-volume) in case something goes south.

1. Pull the latest image 
```bash
docker compose pull
```
2. Recreate the container with the latest image
```bash
docker compose up -d
```

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

## Back up a Volume

It is of vital importance to back up your volume data before updating the software in case something goes south, and 
you lose your database.

### Docker Desktop

If you use Docker Desktop, then 
1. Select the `Volumes` tab to the left
2. Identify the `recipya-data` volume
3. Click the `Export` action button
4. Select `Local file`, select the target directory and click `Export`

### Terminal

Otherwise, run to following command:
```bash
docker run --rm --volumes-from recipya -v $(pwd):/backup ubuntu tar cvf /backup/recipya-volume-backup.tar /root/.config/Recipya
```

## Restore Volume from Backup

Let's say you updated the Docker image and as a result lost your data. Thanks to you having an external backup, the day
will be saved. 

### Docker Desktop

Follow these steps to restore your data using Docker desktop:
1. Select the `Volumes` tab to the left
2. Identify the `recipya-data` volume
3. Click the `Import` action button
4. Select the backed up volume and click `Import`

### Terminal

Otherwise, run the following command:

```bash
docker run --rm --volumes-from recipya -v $(pwd):/backup ubuntu bash -c "cd /root && tar xvf /backup/recipya-volume-backup.tar --strip 1"
```
