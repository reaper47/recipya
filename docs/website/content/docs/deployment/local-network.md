---
title: Local Network
weight: 1
prev: /docs/installation/windows-users
---

The project can be self-hosted network-wide over your local network for access from devices other than the computer
you installed the application on.

## Docker

After [installing the Docker](/docs/installation/docker) image, you can access 
the site network wide via at http://[host computer IP]:[[port]](/docs/installation/config-file).

## Nightly Build

### Linux

Create a service to run the app automatically on boot.

```bash
sudo nano /etc/systemd/system/recipya.service 
```

Copy the following content to the newly-created file.

```bash
[Unit]
Description=Recipya Service
Wants=network.target

[Service]
ExecStart=/path/to/binary/recipya serve

[Install]
WantedBy=multi-user.target
```

Start the service on boot.

```bash
sudo systemctl start recipya.service
sudo systemctl enable recipya.service
```

You can now access the application on your local network at http://[host computer IP]:[[port]](/docs/installation/config-file).

### Windows

The feasibility on Windows remains to be explored.
