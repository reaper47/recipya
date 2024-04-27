---
title: Configuration File
weight: 2
---

The [configuration file](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json)
sets important variables for the application. Let's go over each of them. 

{{< callout type="info" >}}
You don't need to create this file if you don't use Docker because it will be created during the one-time setup.
The admin may change most of these options from the settings.
{{< /callout >}}

- **email**
  - **from**: The email address of your [SendGrid](https://sendgrid.com/) account. Default: `""`.
  - **sendGridAPIKey**: Your [SendGrid](https://app.sendgrid.com/settings/api_keys) API key. The free tier should be sufficient for your needs. Default: `""`.
- **integrations**
  - **azureDocumentIntelligence**
    - **key**: The *KEY 1* variable displayed in the *Keys and endpoint* tab of your [Azure AI Document Intelligence](https://azure.microsoft.com/en-us/products/ai-services/ai-document-intelligence) resource in the [Azure Portal](https://portal.azure.com/#home). Default: `""`.
    - **endpoint**: The *Endpoint* variable displayed in the *Keys and endpoint* tab of your *Document Intelligence* resource in the Azure Portal. Default: `""`.- **server**
  - **autologin**: Whether to login automatically into the application. Useful when you don't need user accounts. Can be `true` or `false`. Default: `false`.
  - **isDemo**: Whether the app is a demo version. Can be `true` or `false`. Default: `false`.
  - **isProduction**: Whether the app is in production. Can be `true` or `false`. Default: `false`.
  - **noSignups**: Whether to disable user account registrations. Set to `true` when you don't want people to create accounts. Default: `false`.
  - **port**: The port the app will be served through if localhost. Is required.
  - **url**: The website the app is served on. This URL will serve as the base link in the emails. Is required.

### Deprecations

The following table lists deprecated environment variables in v1.2.0. They will be removed in v1.3.0.

| Variable                                            | Reason                                                                                                                                                                                                                                                     |
|-----------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **integrations.azureComputerVision.resourceKey**    | Replaced the use of [Azure AI Vision](https://azure.microsoft.com/en-us/products/ai-services/ai-vision) to digitize recipes in favor of [Azure AI Document Intelligence](https://azure.microsoft.com/en-us/products/ai-services/ai-document-intelligence). |
| **integrations.azureComputerVision.visionEndpoint** | Same as above.                                                                                                                                                                                                                                             |                                                                                                                                                                                                                                      |
