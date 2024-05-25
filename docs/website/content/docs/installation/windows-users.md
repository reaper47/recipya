---
title: Windows Users
weight: 4
next: /docs/deployment
---

Please follow these [instructions](/guide/docs/installation/build/) to install Recipya without Docker.

## Network-Wide Access

If you intend to access the website on other devices within your home network, please ensure that the program is permitted through the firewall. To verify:

1. Open the "Firewall & network protection setting"
2. Click on "Allow an app through firewall"
3. Scroll down to recipya*.exe
4. Ensure private and public boxes are checked
5. Apply the settings
6. Find the IP address of your machine (Wi-Fi settings -> Click on the network you are connected to -> IPv4 address)
7. On your other device, access http://[IPv4 address]:[port]

## Windows Defender

If you executed the `recipya-windows-{arch}.exe` binary from the GitHub releases page and Windows Defender quarantined the 
executable, then you should add the folder or file to the exclusion list. This happens because the builds are not signed yet.

1. Unzip a fresh instance of the build.
2. Open Windows Security.
3. Select `Virus & threat protection`.
4. Click `Manage settings` under `Virus & threat protection settings`.
5. Click `Add or remove exclusions` under `Exclusions`.
6. Click the `Add an exclusion` button, select file, and select the executable.
