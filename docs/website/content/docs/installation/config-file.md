---
title: Configuration File
weight: 2
---

The [configuration file](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json)
sets important variables for the application. Let's go over each of them. 

{{< callout type="info" >}}
You do not need to create this file if you do not use Docker because it will be created during the one-time setup.

The admin may change most of these options from the settings.
{{< /callout >}}

- **email**
  - **from**: The administrator's email address. It is usually the email address of your [SendGrid](https://sendgrid.com/) account.
  - **sendGridAPIKey**: Your [SendGrid](https://app.sendgrid.com/settings/api_keys) API key. The free tier should be sufficient for your needs.
- **integrations**
  - **azureComputerVision**
    - **resourceKey**: The *KEY 1* variable displayed on the *Keys and endpoint* tab of your Computer vision resource in the [Azure Portal](https://portal.azure.com/#home).
    - **visionEndpoint**: The *Endpoint* variable displayed on the *Keys and endpoint* tab of your Computer vision resource in the [Azure Portal](https://portal.azure.com/#home).
- **server**
  - **autologin**: Whether to login automatically into the application. Useful when you don't need user accounts. Its value can be either *true* or *false*.
  - **isDemo**: Whether the app is a demo version. Its value can be either *true* or *false*.
  - **isProduction**: Whether the app is in production. Its value can be either *true* or *false*.
  - **noSignups**: Whether to disable user account registrations. Set to *true* when you don't want people to create accounts. Otherwise, *false*.
  - **port**: The port the app will be served through if localhost.
  - **url**: The website the app is served on. This URL will serve as the base link in the emails.

