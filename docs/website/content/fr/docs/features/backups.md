---
title: Sauvegardes
weight: 2
---

Il existe deux types de sauvegardes de données, chacune effectuée une fois tous les trois jours.
Un maximum de dix sauvegardes sont stockées pour chaque type, la sauvegarde la plus ancienne
ayant donc un mois.

### Global

Une sauvegarde globale est une sauvegarde qui enregistre l'état actuel des données de l'application.
Elle est stockée sous `path/to/recipya/data/backup/global/`.

Sa structure est la suivante :

{{< filetree/container >}}
    {{< filetree/folder name="recipya.{year-month-day}.zip" >}}
        {{< filetree/folder name="Recipya" state="open" >}}
            {{< filetree/folder name="Database" state="open" >}}
                {{< filetree/file name="recipya.db" >}}
            {{< /filetree/folder >}}
            {{< filetree/folder name="Images" state="open" >}}
                {{< filetree/folder name="Thumbnails" state="closed" >}}
                    {{< filetree/file name="{uuid_1}.webp" >}}
                    {{< filetree/file name="{uuid_2}.webp" >}}
                    {{< filetree/file name="{uuid_...}.webp" >}}
                    {{< filetree/file name="{uuid_N}.webp" >}}
                {{< /filetree/folder >}}
                {{< filetree/file name="{uuid_1}.webp" >}}
                {{< filetree/file name="{uuid_2}.webp" >}}
                {{< filetree/file name="{uuid_...}.webp" >}}
                {{< filetree/file name="{uuid_N}.webp" >}}
            {{< /filetree/folder >}}
            {{< filetree/folder name="Logs" state="open" >}}
                {{< filetree/file name="recipya.log" >}}
            {{< /filetree/folder >}}
            {{< filetree/folder name="Videos" state="open" >}}
                {{< filetree/file name="{uuid_1}.webm" >}}
                {{< filetree/file name="{uuid_2}.webm" >}}
                {{< filetree/file name="{uuid_...}.webm" >}}
                {{< filetree/file name="{uuid_N}.webm" >}}
            {{< /filetree/folder >}}
            {{< filetree/file name="config.json" >}}
        {{< /filetree/folder >}}
    {{< /filetree/folder >}}
{{< /filetree/container >}}

### Utilisateur

Une sauvegarde utilisateur est une sauvegarde qui enregistre l'état actuel des données d'un utilisateur. Les éléments suivants sont enregistrés :
- Recettes
- Livres de cuisine
- Recettes partagées
- Livres de cuisine partagés

Les sauvegardes utilisateur sont stockées sous `path/to/recipya/data/Backup/users/{userID}`.

Its structure is as follows:
{{< filetree/container >}}
    {{< filetree/folder name="recipya.{year-month-day}.zip" >}}
        {{< filetree/file name="recipes.zip" >}}
        {{< filetree/file name="backup-deletes.sql" >}}
        {{< filetree/file name="backup-inserts.sql" >}}
    {{< /filetree/folder >}}
{{< /filetree/container >}}

## Restaurer

Il est possible de restaurer une sauvegarde précédente. Les instructions pour procéder dépendent de son type.

### Sauvegarde globale

La restauration d'une sauvegarde globale est effectuée uniquement par la personne ayant accès au serveur.

1. Fermez l'application
2. Accédez à `path/to/recipya/data/Backup/global/`
3. Décompressez la sauvegarde que vous souhaitez restaurer
4. Remplacez le contenu sous `path/to/recipya/data/*` par celui de la sauvegarde décompressée
5. Le cas échéant, supprimez `path/to/recipya/data/Database/recipya.db-shm` et `path/to/recipya/data/Database/recipya.db-wal`
6. Démarrez l'application

### Utilisateur

La restauration d'une sauvegarde utilisateur s'effectue via l'interface utilisateur de l'application Web.

1. Accédez à la boîte de dialogue des paramètres
2. Cliquez sur l'onglet « Données »
3. Identifiez le paramètre **Restaurer à partir de la sauvegarde**
4. Sélectionnez la date de sauvegarde
5. Cliquez sur l'icône de lancement de fusée

![](images/settings-restore-backup.webp)
