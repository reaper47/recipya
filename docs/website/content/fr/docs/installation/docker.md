---
title: Docker
weight: 3
---

Une image Docker [stable](https://hub.docker.com/layers/reaper99/recipya/v1.2.0/images/sha256-a32780f33d0c50388ebb955fb837f6ee9b1339987e3898a53ddd5b1f737c7f6e?context=repo) est produite à chaque version.
Une image [instable](https://hub.docker.com/layers/reaper99/recipya/nightly/images/sha256-b2238a11a53982953df5bbcfd7796a19fa382abf75d316b62fa05ac1c867332c?context=repo) est produite le soir à chaque fois la branche `main` a de nouveaux commits.

Vous pouvez installer l'application à l'aide de [Docker](https://www.docker.com/) ou de [Docker Compose](https://docs.docker.com/compose/).

## Utiliser Docker

Il faut d'abord tirer l'image.

```bash
docker pull reaper99/recipya:nightly
```

Ensuite, exécutez l'image. Les variables d'environnement `-e` sont décrites ci-dessous.

```bash
docker run -d \
  --name recipya
  --restart unless-stopped
  -p 8085:8078 \
  -v recipya-data:/root/.config/Recipya \
  -e RECIPYA_SERVER_PORT=8078 \
  reaper99/recipya:nightly
```

Recipya est accessible depuis votre machine hôte à l'adresse [http://localhost:8085](http://localhost:8085).

### Mettre à jour votre conteneur

Exécutez la commande ci-dessous pour une mise à jour rapide avec Watchtower.

N'oubliez pas de remplacer `recipya` par le nom réel de votre conteneur s'il diffère et de [sauvegarder les données de votre volume](#sauvegarder-un-volume).

```bash
docker run --rm --volume /var/run/docker.sock:/var/run/docker.sock containrrr/watchtower --run-once recipya
```

### Synology NAS

La commande Docker suivante devrait être utilisée si vous utilisez Synology.

```bash
docker run -d \
  --name recipya
  --restart unless-stopped
  -p 8085:8078 \
  -v /shared/path/here:/root/.config/Recipya/:rw \
  -e RECIPYA_SERVER_PORT=8078 \
  reaper99/recipya:nightly
```

## Utiliser Docker Composer

Vous pouvez utiliser Docker Compose pour exécuter le conteneur. Tout d'abord, téléchargez le fichier [compose.yaml](https://github.com/reaper47/recipya/blob/main/deploy/compose.yaml).
Modifier les sections `environment` et `ports`. Les variables d'environnement sont décrites ci-dessous. Ensuite, démarrez l'application.

```bash
docker-compose up -d
```

Accédez à l'application via votre navigateur à l'adresse `http://localhost:[host port]`.

Si vous utilisez Windows et que vous avez l'intention d'accéder à l'application sur d'autres appareils de votre réseau domestique, 
assurez-vous d'« Autoriser la connexion » de la règle de pare-feu Windows Defender entrant « Docker Desktop Backend ».

### Mettre à jour votre conteneur

Suivez ces étapes pour mettre à jour Recipya. N'oubliez pas de [sauvegarder vos données de volume](#sauvegarder-un-volume) au cas où quelque chose tournerait mal.

1. Tirez la dernière image
```bash
docker compose pull
```
2. Recréez le conteneur avec la dernière image
```bash
docker compose up -d
```

## Variables d'environnement

| Variable                    | Description                                                                                                                                                                                                                                                         |
|-----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RECIPYA_DI_ENDPOINT         | La variable **Endpoint** affichée dans l'onglet *Clés et point de terminaison* de votre [Azure AI Intelligence documentaire](https://azure.microsoft.com/en-us/products/ai-services/ai-document-intelligence) ressource dans le portail Azure.<br>Par défaut: `""`. |
| RECIPYA_DI_KEY              | La variable **KEY 1** affichée dans l'onglet *Clés et point de terminaison* de votre ressource *Document Intelligence* dans le [portail Azure](https://portal.azure.com/#home).<br>Par défaut: `""`.                                                                |
| RECIPYA_EMAIL               | L'adresse e-mail de votre compte [SendGrid](https://sendgrid.com/).<br>Par défaut: `""`.                                                                                                                                                                            |
| RECIPYA_EMAIL_SENDGRID      | Your [SendGrid](https://app.sendgrid.com/settings/api_keys) API key. The free tier should be sufficient for your needs.<br>Default: `""`.                                                                                                                           |
| RECIPYA_SERVER_AUTOLOGIN    | Si vous souhaitez vous connecter automatiquement à l'application. Peut être `true` ou `false`.<br>Par défaut: `false`.                                                                                                                                              |
| RECIPYA_SERVER_BYPASS_GUIDE | S'il faut accéder directement à la page de connexion en cas d'anonymat. Peut être `true` ou `false`.<br>Par défaut: `false`.                                                                                                                                        |
| RECIPYA_SERVER_IS_DEMO      | Si l'application est une version de démonstration. Peut être `true` ou `false`.<br>Par défaut: `false`.                                                                                                                                                             |
| RECIPYA_SERVER_IS_PROD      | Si l'application est en production. Peut être `true` ou `false`.<br>Par défaut : `false`.                                                                                                                                                                           |
| RECIPYA_SERVER_NO_SIGNUPS   | S'il faut désactiver les enregistrements de comptes utilisateur. Défini sur « true » lorsque vous ne souhaitez pas que les gens créent des comptes.<br>Par défaut: « false ».                                                                                       |
| RECIPYA_SERVER_PORT         | Le port via lequel l'application sera servie si localhost.<br>**Est requis**.                                                                                                                                                                                       |
| RECIPYA_SERVER_URL          | Le site Web sur lequel l'application est diffusée. Cette URL servira de lien de base dans les e-mails.<br>Par défaut: `http://0.0.0.0`.                                                                                                                             |

## Sauvegarder un volume

Il est d'une importance vitale de sauvegarder vos données de volume avant de mettre à jour le logiciel au cas où quelque chose tournerait mal et que vous perdriez votre base de données.

### Docker Desktop

Si vous utilisez Docker Desktop, alors
1. Sélectionnez l'onglet « Volumes » à gauche
2. Identifiez le volume «recipya-data»
3. Cliquez sur le bouton d'action « Exporter »
4. Sélectionnez « Fichier local », sélectionnez le répertoire cible et cliquez sur « Exporter ».

### Terminale

Sinon, exécutez la commande suivante:
```bash
docker run --rm --volumes-from recipya -v $(pwd):/backup ubuntu tar cvf /backup/recipya-volume-backup.tar /root/.config/Recipya
```

## Restaurer le volume à partir d'une sauvegarde

Supposons que vous ayez mis à jour l'image Docker et que vous ayez par conséquent perdu vos données. Grâce à une sauvegarde externe, la journée sera sauvée.

### Docker Desktop

Suivez ces étapes pour restaurer vos données à l'aide de Docker desktop:
1. Sélectionnez l'onglet « Volumes » à gauche
2. Identifiez le volume «recipya-data»
3. Cliquez sur le bouton d'action « Importer »
4. Sélectionnez le volume sauvegardé et cliquez sur « Importer »

### Terminale

Sinon, exécutez la commande suivante :

```bash
docker run --rm --volumes-from recipya -v $(pwd):/backup ubuntu bash -c "cd /root && tar xvf /backup/recipya-volume-backup.tar --strip 1"
```
