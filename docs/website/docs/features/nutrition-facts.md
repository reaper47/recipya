---
sidebar_position: 6
---

# Nutrition Facts

Many recipes lack nutritional information users may be unwilling to calculate themselves. To address this issue,
Recipya can calculate the nutrition facts for you automatically when adding a recipe.

To enable this feature, access your settings via your avatar, click the `Recipes` tab on the left, and 
check the **Calculate nutrition facts** setting's checkbox. 

![img alt](/img/features/settings-nutrition-facts.png)

This setting is initially disabled because adding a recipe will take up to a few additional seconds. This happens
because querying the nutritional database for every ingredient takes time.

Recipya does its best at calculating the nutritional information based on a recipe's ingredients. However, please 
understand that the calculation is more indicative than absolute truth. The information is based on the U.S. Department 
of Agriculture, Agricultural Research Service's [FoodData Central](https://fdc.nal.usda.gov), which is an integrated 
data system that provides expanded nutrient profile data. Please [open an issue](https://github.com/reaper47/recipya/issues/new?assignees=&labels=bug&projects=&template=bug_report.md&title=Problem+with+nutrition+facts)
on GitHub if you ever notice a recipe with nutritional information that seems vastly inaccurate.

:::note

The calculation can only be done with ingredients written in english because the database is in english only.

:::
