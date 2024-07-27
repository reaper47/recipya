---
title: Guide communautaire
weight: 2
---

Bienvenue dans le guide communautaire du projet !

## Canaux de support

- [Espace Matrix](https://app.element.io/#/room/#recipya:matrix.org)
- [GitHub Issues](https://github.com/reaper47/recipya/issues)
- [GitHub Discussions](https://github.com/reaper47/recipya/discussions)

## Façons de contribuer

Recipya est un projet open source collaboratif. Nous accueillons tous ceux qui souhaitent nous aider à faire de ce gestionnaire de
recettes le meilleur possible ! Votre contribution et vos contributions sont essentielles alors que nous travaillons à la création
une solution étonnante de gestion de recettes. Ne pas savoir coder n'est pas une condition pour contribuer !

Les sous-sections suivantes détaillent les différentes manières dont vous pouvez contribuer.

### Développement de fonctionnalités

N'hésitez pas à travailler sur des fonctionnalités qui ne sont pas [assignées](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement+no%3Aassignee).
Je vous encourage également à ouvrir une [demande de fonctionnalité](https://github.com/reaper47/recipya/issues/new?assignees=&labels=enhancement&projects=&template=feature_request.md&title=)
lorsque vous avez des idées qui peuvent améliorer le logiciel.

Vous n'êtes pas obligé d'implémenter les fonctionnalités vous-même si vous ne vous sentez pas à l'aise. Cependant, si vous le souhaitez, la procédure est la suivante.

1. Vérifiez la liste des fonctionnalités actuellement [demandées](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement+no%3Aassignee).
2. Sélectionnez celle sur laquelle vous souhaitez travailler.
3. Indiquez en commentaire que vous souhaitez la corriger ou envoyez-moi un message dans l'espace [Matrix](https://app.element.io/#/room/#recipya:matrix.org)
   afin que je puisse déplacer la tâche vers la colonne « en cours » du [board](https://github.com/users/reaper47/projects/2) et vous l'attribuer.
4. Forkez le dépôt si vous ne l'avez pas encore fait.
5. Implémentez la fonctionnalité et écrivez des tests.
6. Envoyez les modifications à votre fork.
7. Ouvrez un pull request afin que je puisse fusionner votre travail dans « main ».

{{< callout type="info" >}}
Sachez que travailler sur une fonctionnalité sans ouvrir d'abord un problème sur GitHub peut entraîner un rejet si j'estime que
cela ne convient pas à Recipya.
{{< /callout >}}

### Bugs

N'hésitez pas à signaler des bugs lorsque vous en découvrez ! Veuillez vous assurer que le bug que vous avez trouvé n'a pas été signalé dans le 
[GitHub issues](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Abug) before d'en [ouvrir un](https://github.com/reaper47/recipya/issues/new?assignees=&labels=&projects=&template=bug_report.md&title=)
pour réduire le nombre de doublons.

Vous n'êtes pas obligé de corriger le bug vous-même si vous ne vous sentez pas à l'aise. Cependant, si vous le souhaitez, la procédure est la suivante.

1. Vérifiez la liste des bugs actuellement [enregistrés](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Abug).
2. Sélectionnez un problème sur lequel vous souhaitez travailler.
3. Indiquez en commentaire que vous souhaitez le résoudre ou envoyez-moi un message dans l'espace [Matrix](https://app.element.io/#/room/#recipya:matrix.org)
   afin que je puisse déplacer la tâche vers la colonne « en cours » du [board](https://github.com/users/reaper47/projects/2)
   et vous l'assigner.
4. Forkez le dépôt si vous ne l'avez pas encore fait.
5. Corrigez le bug, testez-le correctement et transmettez les modifications à votre fork.
6. Ouvrez un pull request afin que je puisse fusionner votre travail dans « main ».

### Documentation

Ce site Web est la documentation officielle de Recipya. Il est construit à l'aide de [Hextra](https://imfing.github.io/hextra/), qui
est un générateur de site statique très simple à utiliser et à comprendre. Vous n'avez pas besoin d'ouvrir un problème concernant les mises à jour
de la documentation. N'hésitez pas à mettre à jour comme vous le jugez approprié et à ouvrir une demande d'extraction. Vous pouvez nous aider avec
les traductions, l'ajout d'une langue, la correction des fautes de frappe, l'amélioration de la grammaire, l'ajout de sections, la mise à jour des images, le contrôle de version, etc.

Pour développer la documentation localement, vous devez d'abord [fork](https://github.com/reaper47/recipya/fork) le projet.

Ensuite, ouvrez une invite de commande ou un terminal et accédez à `recipya/docs/website`.

```bash
cd path/to/recipya/docs/website
```

Ensuite, diffusez le site Web localement.

```bash
hugo serve
```

Le site Web devrait s'ouvrir automatiquement dans votre navigateur à l'adresse http://localhost:3000.
Vous êtes maintenant libre de modifier le texte et les modifications seront reflétées dans le navigateur lors de l'enregistrement du fichier.

### Aider les autres

C'est toujours un plaisir d'aider quelqu'un qui a besoin d'aide. Veuillez consulter les [canaux d'assistance](#canaux-de-support)
pour savoir où vous pouvez apporter votre aide.
