---
title: Utilisateurs Windows
weight: 4
next: /docs/integrations
---

Veuillez suivre ces [instructions](/guide/fr/docs/installation/build/) pour installer Recipya sans Docker.

## Accès à l'ensemble du réseau

Si vous avez l'intention d'accéder au site Web sur d'autres appareils au sein de votre réseau domestique, veuillez vous assurer que le programme est autorisé via le pare-feu.
Pour vérifier:

1. Ouvrez le « Paramètre de protection pare-feu et réseau »
2. Cliquez sur « Autoriser une application via le pare-feu »
3. Défiler jusqu'à recipya*.exe
4. Assurez-vous que les cases privées et publiques sont cochées
5. Appliquez les paramètres
6. Recherchez l'adresse IP de votre machine (Paramètres Wi-Fi -> Cliquez sur le réseau auquel vous êtes connecté -> Adresse IPv4)
7. Sur votre autre appareil, accédez http://[adresse IPv4]:[port]

## Windows Defender

Si vous avez exécuté le binaire « recipya.exe » à partir de la page des versions de GitHub et que Windows Defender a mis l'exécutable 
en quarantaine, vous devez alors ajouter le dossier ou le fichier à la liste d'exclusion. Cela se produit car les versions ne sont pas encore signées.

1. Décompressez une nouvelle instance de la build.
2. Ouvrez la sécurité Windows.
3. Sélectionnez « Protection contre les virus et les menaces ».
4. Cliquez sur « Gérer les paramètres » sous « Paramètres de protection contre les virus et les menaces ».
5. Cliquez sur « Ajouter ou supprimer des exclusions » sous « Exclusions ».
6. Cliquez sur le bouton « Ajouter une exclusion », sélectionnez le fichier et sélectionnez l'exécutable.
