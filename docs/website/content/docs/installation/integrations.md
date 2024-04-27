---
title: Integrations
weight: 6
---

Recipya uses the following third-party services to enhance the product.

## SendGrid

[SendGrid](https://sendgrid.com) provides a cloud-based service that assists businesses with email delivery.
They offer a [free plan](https://sendgrid.com/en-us/pricing) that allows you to send up to 100 emails per day.

Within Recipya, the email module is used for the following events:
- Send a confirmation email to a user who registered.
- Send a forgot password email

If none of these reasons persuade you to use this service, then leave the `email.from` and `email.sendGridAPIKey` fields
in the [configuration file](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json) empty. No emails
will then be sent.

## Azure AI Document Intelligence

[Azure AI Document Intelligence](https://azure.microsoft.com/en-us/products/ai-services/ai-vision) is an AI service that
applies advanced machine learning to extract text, key-value pairs, tables, and structures from documents automatically 
and accurately. Microsoft offers a [free plan](https://azure.microsoft.com/en-us/pricing/details/ai-document-intelligence/)
(F0) that allows you to perform up to 500 free transactions per month.

Within Recipya, this service is used to [digitize recipes](/guide/docs/features/recipes/add#scan).

If you do not plan on digitizing recipes, then leave the `integrations.azureDocumentIntelligence.key` and
`integrations.azureDocumentIntelligence.endpoint` fields in the [configuration file](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json)
empty. Leave the `RECIPYA_DI_ENDPOINT` and `RECIPYA_DI_KEY` environment variables empty if you use Docker. 
This feature will then be disabled.

Follow these steps to use this integration.
1. Get an Azure subscription. You can [create one for free]((https://azure.microsoft.com/free/cognitive-services/).
2. Add a [Document Intelligence instance](https://portal.azure.com/#create/Microsoft.CognitiveServicesFormRecognizer) in the Azure portal. You can use the free pricing tier (F0) to try the service.
3. After your resource deploys, select *Keys and Endpoint* under *Resource Management* in the sidebar.
4. Copy *KEY 1* to the respective field in Recipya's setting. Alternatively, you may copy it to your configuration file's **integrations.azureDocumentIntelligence.key** field or `RECIPYA_DI_KEY` environment variable if you use Docker.
5. Copy *Endpoint* to the respective field in Recipya's setting. Alternatively, you may copy it to your configuration file's **integrations.azureDocumentIntelligence.endpoint** field or `RECIPYA_DI_ENDPOINT` environment variable if you use Docker.
6. Restart Recipya and test the *Azure AI Document Intelligence* connection from the settings.
