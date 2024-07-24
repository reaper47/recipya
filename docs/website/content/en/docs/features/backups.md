---
title: Backups
weight: 2
---

There are two types of data backups, each done once every three days.
A maximum of ten backups are stored for each type, resulting in the oldest backup 
being one month old.

### Global

A global backup is one which saves the current state of the application data.
It is stored under `path/to/recipya/data/backup/global/`.

Its structure is as follows:

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

### User

A user backup is one which saves the current state of a user's data. The following is saved:
- Recipes 
- Cookbooks
- Shared recipes
- Shared cookbooks

User backups are stored under `path/to/recipya/data/Backup/users/{userID}`.

Its structure is as follows:
{{< filetree/container >}}
    {{< filetree/folder name="recipya.{year-month-day}.zip" >}}
        {{< filetree/file name="recipes.zip" >}}
        {{< filetree/file name="backup-deletes.sql" >}}
        {{< filetree/file name="backup-inserts.sql" >}}
    {{< /filetree/folder >}}
{{< /filetree/container >}}

## Restore

It is possible to restore a previous backup. The instructions on how to do so depends on its type.

### Global

Restoring a global backup is done only by the one who has access to the server.

1. Close the application
2. Navigate to `path/to/recipya/data/Backup/global/`
3. Unzip the backup you wish to restore
4. Replace the content under `path/to/recipya/data/*` with the one from the unzipped backup
5. If applicable, delete `path/to/recipya/data/Database/recipya.db-shm` and `path/to/recipya/data/Database/recipya.db-wal
6. Start the application

### User

Restoring a user backup is done through the web application's user interface.

1. Access the settings dialog
2. Click the `Data` tab
3. Identify the **Restore from backup** setting
4. Select the backup date
5. Click the rocket launch icon

![](images/settings-restore-backup.webp)
