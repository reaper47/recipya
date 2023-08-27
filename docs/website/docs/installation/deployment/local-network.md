---
sidebar_position: 1
---

# Local Network

The project can be self-hosted network-wide over your local network for access from devices other than the computer
you installed the application on. The feasibility on Windows remains to be explored.

## Docker

This subsection will be written once I figure it out.

## Nightly Build

### Linux

Please first ensure that the [url](/docs/installation/config-file) in the configuration file is set to http://0.0.0.0.

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

## Manual Install

The deployment instructions are the same as the [nightly build](#nightly-build) subsection.
