# Recipya

[![Build](https://img.shields.io/github/actions/workflow/status/reaper47/recipya/go.yml?branch=main&logo=Github)](https://github.com/reaper47/recipya/actions/workflows/go.yml)
[![Report](https://goreportcard.com/badge/github.com/reaper47/recipya)](https://goreportcard.com/report/github.com/reaper47/recipya)
[![Contributions](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/reaper47/recipya/issues)

[![GitHub tag](https://img.shields.io/github/v/tag/reaper47/recipya?include_prereleases&label=version)](https://github.com/reaper47/recipya/tags)
[![License](https://img.shields.io/github/license/reaper47/recipya)](./LICENSE)

[Explore the docs](https://recipya.musicavis.ca/docs/) · 
[Demo](https://recipya-app.musicavis.ca)

## Introduction

A clean, simple and powerful recipe manager web application for unforgettable family recipes, empowering you to curate and share your favorite recipes.
It is focused on simplicity for the whole family to enjoy.

![Recipe page screenshot](.github/screenshot-recipes.webp)

## Features

- Manage your favorite recipes
- Import recipes from around the web
- Digitize paper recipes
- Organize your recipes into cookbooks
- Easily migrate your recipes from [Mealie](https://mealie.io), [Tandoor](https://tandoor.dev) and [Nextcloud Cookbook](https://apps.nextcloud.com/apps/cookbook)
- Automatic conversion to your preferred measurement system (imperial/metric)
- Calculate nutritional information automatically
- Print any recipe in your collection
- Prevent your device from going to sleep while viewing a recipe
- Follows your system's theme (light/dark) or choose among 32 themes
- Cross-compiled for Windows, Linux, and macOS

## Getting Started

### Installation

The installation instructions are written in the [installation section](https://recipya.musicavis.ca/docs/installation/) of the documentation.

### Building the Project

Follow these steps to build the project yourself:
1. Clone the project.
   ```bash
   git clone https://github.com/reaper47/recipya.git
   ```
2. Install the required [dependencies](https://recipya.musicavis.ca/docs/development/intro/#dependencies).
3. [Build](https://recipya.musicavis.ca/docs/development/build/) the project.

Alternatively, you may use the [development container](https://recipya.musicavis.ca/docs/development/devcontainer/).
Recipya's Docker [container](https://github.com/reaper47/recipya/tree/main/.devcontainer) includes all the necessary tools and dependencies 
you need to start writing code quickly.

### Updating

#### Docker

Please follow these [instructions](https://recipya.musicavis.ca/docs/installation/docker/#updating-your-container-1) 
to update your Docker instance.

#### Release build

If you installed a release build of Recipya, i.e. v1.0.0, then the software will notify you once an update is available. 
You can [self-update](https://recipya.musicavis.ca/docs/features/updater) the application from the settings dialog.

Let's explain the mechanism developers. Basically, a cron job is run [every three days](https://github.com/reaper47/recipya/blob/main/internal/jobs/jobs.go#L70-L82)
to check whether there is a new GitHub release. If so, the update indicators in the UI will be enabled. Once the user presses 
the "Update" button, the latest release is fetched, unpacked and the application restarted.

## Contributing

Contributions are always welcome! Please open an issue, start a [discussion](https://github.com/reaper47/recipya/discussions), open a pull request or send an email 
at macpoule@gmail.com. The same applies if you have any feedback or need support.

You can also join our development and support channel on the [Matrix space: #recipya:matrix.org](https://app.element.io/#/room/#recipya:matrix.org).
Matrix is similar to Discord but is open source.

## Sponsors

I am grateful for any support that helps me continue to develop this project and to host it reliably. Your sponsorship will 
help me pay for the SendGrid Essentials plan to increase the number of emails that can be sent. The free plan currently 
used allows sending up to 100 emails per day.

You can sponsor me on 
[GitHub Sponsors](https://github.com/sponsors/reaper47) or
[Buy Me a Coffee](https://www.buymeacoffee.com/macpoule).

Your support is greatly appreciated! A third of donations will be sent to the Armed Forces of Ukraine 🇺🇦

This project is supported by these kind people:
<img src="scripts/sponsors/sponsors.svg" style="width:100%;max-width:800px;"/>

## Other Recipe Manager Apps

- [Tandoor](https://github.com/TandoorRecipes/recipes)
- [Mealie](https://github.com/mealie-recipes/mealie)
- [Paprika](https://www.paprikaapp.com/)
- [Grocy](https://grocy.info/)
- [Cooklist](https://cooklist.com/)
- [Grossr](https://grossr.com/)
- [Awesome List](https://github.com/awesome-selfhosted/awesome-selfhosted#recipe-management)

# Inspiration

This project was mainly coded to blasting the following albums:
- [4am](https://www.youtube.com/watch?v=tBcPji_jRDc)
- [Abysmal Dawn - Phylogenesis](https://www.youtube.com/watch?v=xJMybqRMedk&pp=ygUMYWJ5c21hbCBkYXdu)
- [Archspire - Bleed the Future](https://www.youtube.com/watch?v=o8H9ahswldM)
- [Atavistia - Cosmic Warfare](https://www.youtube.com/watch?v=VjJ_zb4RF2E)
- [Beast In Black - Dark Connection](https://www.youtube.com/watch?v=7NyON-NzBr4)
- [Cattle Decapitation - Terrasite](https://www.youtube.com/watch?v=x6rEDMqM36I)
- [Ensiferum - From Afar](https://www.youtube.com/watch?v=6r8OPu3SRSM)
- [Lofi Girl - lofi hip hop radio](https://www.youtube.com/watch?v=jfKfPfyJRdk)
- [Lofi Girl - synthwave radio](https://www.youtube.com/watch?v=4xDzrJKXOOY)
- [Mozart - Requiem Dm](https://www.youtube.com/watch?v=pBGVfwOLU1w0)
- [Necrophobic - In the Twilight Grey](https://www.youtube.com/watch?v=eDFD6YnMid8)
- [Pain - You Only Live Twice](https://www.youtube.com/watch?v=obgCEoLzLs4)
- [Sonata Arctica - Talviyö](https://www.youtube.com/watch?v=x6rEDMqM36I)
- [Wintersun - Wintersun](https://www.youtube.com/watch?v=W0M3HAMus7g&pp=ygUPd2ludGVyc3VuIGFsYnVt)
