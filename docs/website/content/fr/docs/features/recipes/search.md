---
title: Recherche
weight: 5
---

Vous pouvez rechercher des recettes à partir de la zone de texte de recherche située sous le bouton Ajouter des recettes de la page d'affichage des recettes.
Tapez vos mots-clés délimités par un espace, une virgule ou les deux. Appuyez ensuite sur la touche Entrée pour lancer la recherche.
![](images/recipes-search-view.webp)

## Recherche de base

La recherche de base est celle dans laquelle l'utilisateur saisit quelque chose dans la zone de recherche sans utiliser aucun
des concepts de recherche avancée.

Le mode de recherche par défaut est __Recherche complète__. En d'autres termes, votre requête recherchera les champs suivants de toutes les
recettes vous appartenant :
- Titre
- Description
- Catégorie
- Ingrédients
- Instructions
- Mots-clés
- Source
- Outils

Les résultats seront classés en fonction de leur pertinence par rapport aux termes de recherche.

![](images/recipes-search-query.webp)

## Recherche Avancée

La fonction de recherche avancée vous permet d'adapter la requête de recherche à des besoins spécifiques. Elle est similaire à l'utilisation des fonctions de recherche
avancées de Google. Par exemple, la recherche de `magnetic declination site:.edu` dans Google donnera des résultats
contenant les termes « déclinaison magnétique » pour les sites Web du domaine de premier niveau `.edu`.

Le tableau suivant fournit des exemples de la manière d'effectuer différentes recherches. Vous pouvez combiner n'importe laquelle des recherches dans n'importe quel ordre.

| Recherche                              | Exemple                                                |
|----------------------------------------|--------------------------------------------------------|
| N'importe quel champ                   | big green squash                                       |
| Par catégorie                          | cat:dinner                                             |
| Plusieurs catégories                   | cat:breakfast,dinner                                   |
| Sous-catégorie                         | cat:beverages:cocktails                                |
| N'importe quel champ de catégorie      | chicken cat:dinner                                     |
| Par nom                                | name:chicken kyiv                                      |
| Par nom et par catégorie               | name:chicken kyiv cat:lunch                            |
| N'importe quel champ, nom et catégorie | best name:chicken kyiv cat:lunch                       |
| Par description                        | desc:tender savory stacked                             |
| Plusieurs descriptions                 | desc:tender savory stacked,juicy crispy pieces chicken |
| Par cuisine                            | cuisine:ukrainian                                      |
| Plusieurs cuisines                     | cuisine:ukrainian,japanese                             |
| Par ingrédient                         | ing:onions                                             |
| Plusieurs ingrédients                  | ing:olive oil,thyme,butter                             |
| Par instruction                        | ins:preheat oven 350                                   |
| Plusieurs instructions                 | ins:preheat oven 350,melt butter                       |
| Par mot-clé                            | tag:biscuits                                           |
| Plusieurs mots-clés                    | tag:biscuits,mardi gras                                |
| Par outil                              | tool:frying pan                                        |
| Plusieurs outils                       | tool:frying pan,wok                                    |
| Par source                             | src:allrecipes.com                                     |
| Plusieurs sources                      | src:allrecipes.com,tasteofhome.com                     |

### Aide

Vous pouvez accéder à la boîte de dialogue d'aide à la recherche avancée en cliquant sur le bouton d'information à l'extrême gauche de la barre de recherche.

![](images/recipes-search-help.webp)
