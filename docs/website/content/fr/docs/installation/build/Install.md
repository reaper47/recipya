---
title: Install
weight: 1
prev: /docs/installation/build
---

Cette section cible une installation sans Docker.

Les versions portables et autonomes sont disponibles sur la [page des versions](https://github.com/reaper47/recipya/releases) sur GitHub.
La [version nocturne](https://github.com/reaper47/recipya/releases/tag/nightly) est mise à jour tous les soirs si la branche principale a de nouveaux commits.

Tout d’abord, téléchargez la version du logiciel que vous souhaitez installer compatible avec votre système et extrayez le fichier zip.
Veuillez consulter le tableau des [plateformes prises en charge](/guide/fr/docs/installation/system-requirements) si vous ne savez pas quel fichier télécharger.

Ensuite, démarrez le serveur en ouvrant une invite de commande dans le dossier et exécutez la commande suivante.
L'application effectuera une configuration de démarrage unique si ce n'est déjà fait.

```bash
./recipya serve
```

Vous pouvez désormais accéder au site Internet à l'adresse indiquée.

## Example

Supposons que vous ayez un ordinateur Windows 11 et que vous souhaitiez y installer Recipya v1.2.0.

{{% steps %}}

### Accès

Vous accéderez d'abord à la [page des versions](https://github.com/reaper47/recipya/releases).

### Télécharger

Identifiez la version « v1.2.0 » et téléchargez « recipya-windows-amd64.zip » dans la section **Assets**.

### Extraire

Extrayez le fichier zip sur votre ordinateur.

### Exécuter

130 / 5,000
Ouvrez une [invite de commande](https://en.wikiversity.org/wiki/Command_Prompt/Open) et accédez au dossier que vous avez précédemment extrait.

```text
cd C:\path\to\recipya
```

Ensuite, exécutez Recipya une fois que votre invite de commande est dans son dossier.

```text
.\recipya serve
```

### Profite bien

Ouvrez votre navigateur à l'adresse spécifiée dans le texte de l'invite de commande.

Si vous voyez le résultat suivant:
```text
OK FDC database
OK Configuration file
Recipya is properly set up
File locations:
        - Backups: C:\Users\<user>\AppData\Roaming\Recipya\Backup
        - Database: C:\Users\<user>\AppData\Roaming\Recipya\Database
        - Images: C:\Users\<user>\AppData\Roaming\Recipya\Images
        - Logs: C:\Users\<user>\AppData\Roaming\Recipya\Logs
2024/05/23 07:46:46 goose: no migrations to run. current version: 20240522133726
```

Tout est beau et vous accéderez ensuite à « http://127.0.0.1:8078 ».

{{% /steps %}}
