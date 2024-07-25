---
title: Réseau local
weight: 1
prev: /docs/installation/windows-users
---

Le projet est auto-hébergeable sur l'ensemble de votre réseau local pour y accéder à partir d'appareils autres que
l'ordinateur sur lequel vous l'avez installé.

## Docker

Après avoir [installé l'image Docker](/guide/docs/installation/docker), vous pouvez accéder à l'ensemble du réseau du site à l'adresse http://<host computer IP>:[<port>](/guide/docs/installation/build/config-file/].

## Release Build

### Linux

Créez un service pour exécuter l'application automatiquement au démarrage.

```bash
sudo nano /etc/systemd/system/recipya.service 
```

Copiez le contenu suivant dans le fichier nouvellement créé.

```bash
[Unit]
Description=Recipya Service
Wants=network.target

[Service]
ExecStart=/path/to/binary/recipya serve
Environment=HOME=/root

[Install]
WantedBy=multi-user.target
```

Démarrez le service au démarrage.

```bash
sudo systemctl start recipya.service
sudo systemctl enable recipya.service
```

Vous pouvez désormais accéder l'application sur votre réseau local à l'adresse http://<host computer IP>:[<port>](/docs/installation/build/config-file)].

### Windows

La faisabilité sur Windows reste à explorer.
