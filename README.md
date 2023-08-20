# Recipya

A beautiful recipe manager web application for unforgettable family recipes, empowering you to curate and share your favorite recipes. It is focussed on simplicity for the whole family to enjoy.

## Features

- Manage your favorite recipes
- Import recipes from around the web
- Digitize paper recipes (planned)
- Works seamlessly with recipes in your [Nextcloud Cookbook](https://apps.nextcloud.com/apps/cookbook)
- Automatic conversion to your preferred measurement system (imperial/metric)
- Print any recipe in your collection
- Follows your system's theme (light/dark)
- Cross-compiled for Windows, Linux, and macOS

## Screenshots

Screenshots will be added later, once the project is more mature.

## Demo

The demo link will be added later, once I host the app somewhere.

## Documentation

The documentation has to be written.

## Installation

### Nightly Build

A portable, standalone build is available on the [releases page](https://github.com/reaper47/recipya/releases/tag/nightly). It is updated nightly.

The following table lists the supported platforms and devices. It will help you decide which zip file to download if you do not know your computer's architecture.

| Platform      | Explanation                                  | Device Examples                                         |
|---------------|----------------------------------------------|---------------------------------------------------------|
| darwin/amd64  | macOS on 64-bit Intel (x86-64) architecture  | Apple MacBook, iMac, Mac Mini, Mac Pro                  |
| darwin/arm64  | macOS on ARM64 architecture                  | MacBook Air (M1), MacBook Pro (M1), Mac Mini (M1)       |
| linux/386     | Linux on 32-bit x86 architecture             | Older PCs, embedded systems                             |
| linux/amd64   | Linux on 64-bit x86 architecture (x86-64)    | Desktops, laptops, servers, cloud instances             |
| linux/arm     | Linux on ARMv6 architecture                  | Raspberry Pi 1st gen, IoT devices, some old smartphones |
| linux/arm64   | Linux on ARMv8 64-bit architecture           | Raspberry Pi 3rd/4th gen, modern smartphones            |
| linux/riscv64 | Linux on 64-bit RISC-V architecture          | RISC-V development boards, experimental devices         |
| linux/s390x   | Linux on IBM System z architecture           | IBM mainframes, servers                                 |
| windows/amd64 | Windows on 64-bit x86 architecture           | Modern Windows PCs, servers, virtual machines           |
| windows/arm64 | Windows on ARM64 architecture                | Microsoft Surface Pro X, ARM-based Windows devices      |

The program cannot be updated within the interface yet. If you wish to update the app, please download the nightly build and replace your existing `recipya.exe` with the one from the build zip.

Extract the zip once downloaded. Then, follow these [steps](#modify-configuration-file) to modify your configuration file.

Finally, start the server by opening a command prompt in the folder, and run the following command.

```bash
./recipya.exe serve
```

You can now access the website at the address specified.

#### Note for Windows users

If you intend to access the website on other devices within your home network, please ensure that the program is permitted through the firewall. To verify:
1. Open the "Firewall & network protection setting"
1. Click on "Allow an app through firewall"
1. Scroll down to recipya*.exe
1. Make sure private and public boxes are checked
1. Apply the settings
1. Find the IP address of your machine (Wi-Fi settings -> Click on the network you are connected to -> IPv4 address)
1. On your other device, access http://[IPv4 address]:[port]

### Docker

A Docker image called `reaper99/recipya` is produced nightly.

#### Using Docker

You first have to fetch it.

```bash
docker pull reaper99/recipya:nightly
```

Then, run the image. You must pass your `config.json` file to the container.

```bash
docker run -v path/to/config.json:/app/config.json -p [host port]:[port specified in config.json] -d reaper99/recipya:nightly reaper99/recipya:nightly
```

#### Using Docker Compose

You can use Docker Compose to run the container. First, you need to modify the ports
and the path to your local config.json in the [compose.yaml](https://github.com/reaper47/recipya/blob/main/containers/compose.yaml).
Then, start the application.

```bash
docker-compose up -d
```

Access the app through your browser at `http://localhost:[host port]`.

If you are using Windows and you intend to access the app on other devices within your home network, please ensure
to `Allow the connection` of the `Docker Desktop Backend` inbound Windows Defender Firewall rule.

### Manual Install

#### Windows

The `make` program is required to build the project. To verify whether it's installed on your machine, execute `make` in a command prompt or PowerShell.

Follow these steps if not installed:

1. Open either the Command Prompt or Powershell in administrator mode.
1. Execute `winget install GnuWin32.Make`
1. Add `C:\Program Files (x86)\GnuWin32\bin` to the Windows PATH environment variable.

---

Follow these steps to build the program:

First, clone the repository.

```bash
git clone https://github.com/reaper47/recipya.git
```

Go to the project directory.

```bash
cd recipya
```

Build the project.

```bash
make
```

Follow these [steps](#modify-configuration-file) to modify your configuration file.

Start the server.

```bash
.\bin\recipya serve
```

Finally, please ensure that the program is permitted through the firewall if you intend to access it on other devices within your home network.

#### Linux and macOS

There is no nightly build available for Linux and macOS. You will have to build it yourself.

First, clone the repository.

```bash
git clone https://github.com/reaper47/recipya.git
```

Go to the project directory.

```bash
cd recipya
```

Build the project.

```bash
make
```

Follow these [steps](#modify-configuration-file) to modify your configuration file.

Start the server.

```bash
./bin/recipya serve
```

### Modify Configuration File

Once the build zip is extracted or the project built, rename *config.json.example* to *config.json* and open the file to edit the following variables:
- **email.from**: The administrator's email address
- **email.sendGridAPIKey**: Your [SendGrid](https://sendgrid.com/) API key. The free tier should be sufficient for your needs.
- **port**: The port the app will be served through.
- **url**: The website the app is served on. This URL is used in the emails.


## Deployment

The project can be self-hosted. You first need to înstall the project [manually](#manual-install).

Then, create a service to run the app automatically on boot.

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
## Running Tests

Execute the following command to run the tests.

```bash
make test
```

## Contributing

Contributions are always welcome! Please open an issue, start a [discussion](https://github.com/reaper47/recipya/discussions), open a pull request or send an email 
at macpoule@gmail.com. The same applies if you have any feedback or need support. 

You can also join our development and support channel on 
the [Matrix space: #recipya:matrix.org](https://app.element.io/#/room/#recipya:matrix.org).
Matrix is similar to Discord but is open source.

## Tech Stack

**Client:** HTML, htmx, _hyperscript, TailwindCSS

**Server:** Go

## Other Recipe Manager Apps

- [Mealie](https://github.com/mealie-recipes/mealie)
- [Paprika](https://www.paprikaapp.com/)
- [Grocy](https://grocy.info/)
- [Cooklist](https://cooklist.com/)
- [Grossr](https://grossr.com/)
- [Awesome List](
https://github.com/awesome-selfhosted/awesome-selfhosted#recipe-management)

# Inspiration

This project was mainly coded to blasting the following albums:
- [Archspire - Bleed the Future](https://www.youtube.com/watch?v=o8H9ahswldM)
- [Sonata Arctica - Talviyö](https://www.youtube.com/watch?v=x6rEDMqM36I)
- [Cattle Decapitation - Terrasite](https://www.youtube.com/watch?v=x6rEDMqM36I)
- [Mozart - Requiem Dm](https://www.youtube.com/watch?v=pBGVfwOLU1w0)
- [Ensiferum - From Afar](https://www.youtube.com/watch?v=6r8OPu3SRSM)
- [Pain - You Only Live Twice](https://www.youtube.com/watch?v=obgCEoLzLs4)
- [Abysmal Dawn - Phylogenesis](https://www.youtube.com/watch?v=xJMybqRMedk&pp=ygUMYWJ5c21hbCBkYXdu)
- [Lofi Hip Hop Radio](https://www.youtube.com/watch?v=jfKfPfyJRdk)
- [4am](https://www.youtube.com/watch?v=tBcPji_jRDc)
