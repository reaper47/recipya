---
title: Prérequis
weight: 1
next: /docs/installation/build/
prev: /docs/installation/
---

## Configuration requise

Le tableau suivant répertorie les plateformes et appareils pris en charge.

| Plateforme    | Explication                                   | Exemples d'appareils                                                    |
|---------------|-----------------------------------------------|-------------------------------------------------------------------------|
| darwin/amd64  | macOS sur architecture Intel 64 bits (x86-64) | Apple MacBook, iMac, Mac Mini, Mac Pro                                  |
| darwin/arm64  | macOS sur architecture ARM64                  | MacBook Air (M1), MacBook Pro (M1), Mac Mini (M1)                       |
| linux/amd64   | Linux sur architecture x86 64 bits (x86-64)   | Ordinateurs de bureau, ordinateurs portables, serveurs, instances cloud |
| linux/arm64   | Linux sur architecture ARMv8 64 bits          | Raspberry Pi 3e/4e génération, smartphones modernes                     |
| windows/amd64 | Windows sur architecture x86 64 bits          | PC, serveurs et machines virtuelles Windows modernes                    |
| windows/arm64 | Windows sur architecture ARM64                | Microsoft Surface Pro X, appareils Windows basés sur ARM                |

De plus, vous devez disposer d'au moins 300 Mo d'espace libre.

## Compatibilité du navigateur

| Navigateur | Version |     Compatibilité      |
|------------|:-------:|:----------------------:|
| Brave      |   37+   |  {{< icon "check" >}}  |
| Chrome     |  114+   |  {{< icon "check" >}}  |
| Edge       |  114+   |  {{< icon "check" >}}  |
| Firefox    |  125+   |  {{< icon "check" >}}  |
| IE         |   N/A   |    {{< icon "x" >}}    |
| Opera      |  100+   |  {{< icon "check" >}}  |
| Safari     |   17+   |  {{< icon "check" >}}  |
| Vanadium   |  114+   |  {{< icon "check" >}}  |

## Dépendances logicielles

| Logiciel | Version |      Facultatif      |
|----------|:-------:|:--------------------:|
| FFmpeg   |   7+    | {{< icon "check" >}} |

### FFmpeg

[FFmpeg](https://en.wikipedia.org/wiki/FFmpeg) est utilisé pour convertir des fichiers vidéo au format de fichier multimédia audiovisuel [WebM](https://en.wikipedia.org/wiki/WebM).
Il est inclus dans l'image Docker. Sinon, il sera automatiquement installé si vous utilisez Windows.
Si vous utilisez macOS ou Linux, vous devez l'installer manuellement.

La fonction vidéo sera désactivée si FFmpeg n'est pas installé. Vous n'avez pas besoin de ce logiciel si vous n'avez pas l'intention d'ajouter des vidéos aux recettes.
