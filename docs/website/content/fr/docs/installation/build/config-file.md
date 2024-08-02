---
title: Configuration File
weight: 2
next: /docs/installation/docker
---

Le [fichier de configuration](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json) définit des variables importantes pour l'application. 
Examinons chacun d'eux.

{{< callout type="info" >}}
L'administrateur peut modifier la plupart de ces options à partir des paramètres.
{{< /callout >}}

- **email**
  - **from**: L'adresse e-mail de votre compte [SendGrid](https://sendgrid.com/). Défaut: `""`.
  - **sendGridAPIKey**: Votre clé API [SendGrid](https://app.sendgrid.com/settings/api_keys). Le niveau gratuit devrait être suffisant pour vos besoins. Défaut: `""`.
- **integrations**
  - **azureDocumentIntelligence**
    - **key**: La variable *CLÉ 1* affiché dans l'onglet *Clés et point de terminaison* de votre ressource [Azure AI Document Intelligence](https://azure.microsoft.com/en-us/products/ai-services/ai-document-intelligence) dans le [portail Azure](https://portal.azure.com/#home). Défaut:      `""`.
    - **endpoint**: La variable *Endpoint* affichée dans l'onglet *Clés et point de terminaison* de votre ressource *Document Intelligence* dans le portail Azure. Défaut: `""`.
- **server** 
  - **autologin**: S'il faut se connecter automatiquement à l'application. Utile lorsque vous n'avez pas besoin de comptes d'utilisateurs. Peut être `true` ou `false`. Défault: `false`.
  - **isDemo**: Si l’application est une version démo. Peut être `true` ou `false`. Défault: `false`.
  - **isProduction**: Si l'application est en production. Peut être `true` ou `false`. Défault: `false`.
  - **noSignups**: S'il faut désactiver les enregistrements de comptes utilisateur. Définissez sur `true` lorsque vous ne souhaitez pas que les gens créent des comptes. Défault: `false`.
  - **port**: Le port via lequel l'application sera servie s'il est localhost. __Est requis__.
  - **url**: Le site web sur lequel l'application est diffusée. Cette URL servira de lien de base dans les e-mails. Défault: `http://0.0.0.0`.
