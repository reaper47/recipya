---
title: Prerequisites
weight: 1
next: /docs/installation/build/
prev: /docs/installation/
---

## System Requirements

The following table lists the supported platforms and devices.

| Platform      | Explanation                                  | Device Examples                                         |
|---------------|----------------------------------------------|---------------------------------------------------------|
| darwin/amd64  | macOS on 64-bit Intel (x86-64) architecture  | Apple MacBook, iMac, Mac Mini, Mac Pro                  |
| darwin/arm64  | macOS on ARM64 architecture                  | MacBook Air (M1), MacBook Pro (M1), Mac Mini (M1)       |
| linux/amd64   | Linux on 64-bit x86 architecture (x86-64)    | Desktops, laptops, servers, cloud instances             |
| linux/arm64   | Linux on ARMv8 64-bit architecture           | Raspberry Pi 3rd/4th gen, modern smartphones            |
| windows/amd64 | Windows on 64-bit x86 architecture           | Modern Windows PCs, servers, virtual machines           |
| windows/arm64 | Windows on ARM64 architecture                | Microsoft Surface Pro X, ARM-based Windows devices      |

In addition, you must have at least 300 MB of free space.

## Browser Compatibility

| Browser  | Version |     Compatibility      |
|----------|:-------:|:----------------------:|
| Brave    |   37+   |  {{< icon "check" >}}  |
| Chrome   |  114+   |  {{< icon "check" >}}  |
| Edge     |  114+   |  {{< icon "check" >}}  |
| Firefox  |  125+   |  {{< icon "check" >}}  |
| IE       |   N/A   |    {{< icon "x" >}}    |
| Opera    |  100+   |  {{< icon "check" >}}  |
| Safari   |   17+   |  {{< icon "check" >}}  |
| Vanadium |  114+   |  {{< icon "check" >}}  |

## Dependencies

| Software | Version |       Optional       |
|----------|:-------:|:--------------------:|
| FFmpeg   |   7+    | {{< icon "check" >}} |

### FFmpeg

[FFmpeg](https://en.wikipedia.org/wiki/FFmpeg) is used to convert video files to the [WebM](https://en.wikipedia.org/wiki/WebM) 
audiovisual media file format. It is included in the Docker image. Otherwise, it will be automatically installed if you use Windows.
If you use macOS or Linux, then you should install it manually.

The video feature will be disabled if FFmpeg is not installed.
