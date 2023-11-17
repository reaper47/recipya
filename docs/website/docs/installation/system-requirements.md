---
sidebar_position: 1
---

# Prerequisites

## System Requirements

The following table lists the supported platforms and devices.

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

In addition, you must have at least 300 MB of free space.

## Third-party Services

Recipya uses the following third-party services to enhance the product.

### SendGrid

[SendGrid](https://sendgrid.com) provides a cloud-based service that assists businesses with email delivery.
They offer a [free plan](https://sendgrid.com/en-us/pricing) that allows you to send up to 100 emails per day.

Within Recipya, the email module is used for the following events:
- Send a confirmation email to a user who registered.
- Notifying the administrator when some errors occur.
- Inform the administrator when a user requests a URL to import from an unsupported website.

If none of these reasons persuade you to use this service, then leave the `email.from` and `email.sendGridAPIKey` fields
in the [configuration file](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json) empty. No emails
will then be sent.

### Azure AI Vision

[Azure AI Vision](https://azure.microsoft.com/en-us/products/ai-services/ai-vision) is a unified service that offers 
innovative computer vision capabilities. It gives apps the ability to analyze images, read text, and detect faces 
with prebuilt image tagging, text extraction with optical character recognition (OCR), and responsible facial 
recognition. Microsoft offers a [free plan](https://azure.microsoft.com/en-us/pricing/details/cognitive-services/computer-vision/)
that allows you to perform 5000 transactions per month.

Within Recipya, this service is used to [digitize recipes](/features/recipes/add#scan).

If you do not plan on digitizing paper recipes, then leave the `integrations.azureComputerVision.resourceKey` and
`integrations.azureComputerVision.visionEndpoint` fields in the [configuration file](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json) 
empty. This feature will then be disabled.
