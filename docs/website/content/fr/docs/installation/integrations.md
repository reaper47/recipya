---
title: Intégrations
weight: 5
---

Recipya utilise les services tiers suivants pour améliorer l'expérience du produit.

## SendGrid

[SendGrid](https://sendgrid.com) fournit un service cloud qui aide les entreprises à envoyer des e-mails.
Ils offrent un [forfait gratuit](https://sendgrid.com/en-us/pricing) qui vous permet d'envoyer jusqu'à 100 e-mails par jour.

Au sein de Recipya, le module email est utilisé pour les événements suivants :
- Envoyer un email de confirmation à un utilisateur inscrit.
- Envoyer un e-mail de mot de passe oublié

Si aucune de ces raisons ne vous persuade d'utiliser ce service, laissez les champs `email.from` et `email.sendGridAPIKey` dans le 
[fichier de configuration](https://github.com/reaper47/recipya/blob/main/) déployer/config.example.json) vide. Aucun email ne sera alors envoyé.

## Azure AI Intelligence documentaire

[Azure AI Intelligence documentaire](https://azure.microsoft.com/en-us/products/ai-services/ai-vision) est un service d'intelligence 
artificielle qui applique l'apprentissage automatique avancé pour extraire automatiquement et avec précision du texte, des paires clé-valeur, 
des tableaux et des structures à partir de documents. Microsoft offre un [plan gratuit](https://azure.microsoft.com/en-us/pricing/details/ai-document-intelligence/)
(F0) qui vous permet d'effectuer jusqu'à 500 transactions gratuites par mois.

Au sein de Recipya, ce service est utilisé pour [numériser des recettes](/guide/fr/docs/features/recipes/add#scan).

Si vous ne prévoyez pas de numériser des recettes, laissez les champs `integrations.azureDocumentIntelligence.key` et `integrations.azureDocumentIntelligence.endpoint`
dans le [fichier de configuration](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json) vide.
Laissez les variables d'environnement `RECIPYA_DI_ENDPOINT` et `RECIPYA_DI_KEY` vides si vous utilisez Docker.
Cette fonctionnalité sera alors désactivée.

Suivez ces étapes pour utiliser cette intégration:
1. Obtenez un abonnement Azure. Vous pouvez [en créer un gratuitement](https://azure.microsoft.com/free/cognitive-services/).
2. Ajoutez une [instance Intelligence documentaire](https://portal.azure.com/#create/Microsoft.CognitiveServicesFormRecognizer) dans le portail Azure. Vous pouvez utiliser le niveau tarifaire gratuit (F0) pour essayer le service.
3. Sous **Détails de l'instance**, sélectionnez **Région** _East US_, _West US2_ ou _West Europe_. Les autres régions sont incompatibles avec cette ressource.
4. Une fois votre ressource déployée, sélectionnez *Clés et point de terminaison* sous *Gestion des ressources* dans la barre de navigation.
   ![alt text](https://learn.microsoft.com/en-us/azure/ai-services/document-intelligence/media/containers/keys-and-endpoint.png?view=doc-intel-3.1.0)
5. Copiez *KEY 1* dans le champ correspondant dans les paramètres de Recipya. Vous pouvez également le copier dans le champ **integrations.azureDocumentIntelligence.key** de votre fichier de configuration ou dans la variable d'environnement `RECIPYA_DI_KEY` si vous utilisez Docker.
6. Copiez *Endpoint* dans le champ correspondant dans les paramètres de Recipya. Vous pouvez également le copier dans le champ **integrations.azureDocumentIntelligence.endpoint** de votre fichier de configuration ou dans la variable d'environnement `RECIPYA_DI_ENDPOINT` si vous utilisez Docker.
7. Redémarrez Recipya et testez la connexion *Azure AI Document Intelligence* à partir des paramètres.

### Limitations

- Pour les PDF et TIFF, jusqu'à 2000 pages peuvent être traitées (avec un abonnement gratuit, seules les deux premières pages sont traitées).
- La taille du fichier pour l'analyse des documents est de 500 Mo pour le niveau payant (S0) et de 4 Mo pour le niveau gratuit (F0).
- Si vos fichiers PDF sont verrouillés par mot de passe, vous devez supprimer le verrou avant de les soumettre.
