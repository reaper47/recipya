---
sidebar_position: 1
---

# Add

You can add a recipe by clicking the blue **Add Recipe** in the middle of the application bar.

![img alt](/img/features/add-recipe.png)

## Adding Recipes

You will be presented with four different ways of adding a recipe to your collection.

- [Manual](#manual)
- [Scan](#scan)
- [Website](#website)
- [Import](#import)

![img alt](/img/features/add-recipe-options.png)

### Manual

The simplest method of inputting a recipe involves completing a form that outlines your dish.
Mandatory fields are indicated with an asterisk (*).

![img alt](/img/features/add-recipe-manual.png)

You might find these shortcuts useful if you are filling the form from your computer.

| Shortcut     | How to enable                     | Result              |
|--------------|-----------------------------------|---------------------|
| Enter        | Focus on an ingredient text input | Add ingredient row  |
| Ctrl + Enter | Focus on an instruction text area | Add instruction row |

You can also reorder the ingredients and the instructions by dragging the arrow vertically.

### Scan

You can upload an image or a scan of a handwritten or printed recipe to add it to your collection. 
This option is useful for digitizing your and your family's paper recipes.

:::danger

This feature is not implemented yet.

:::

### Website

You can import any recipe from the supported websites. To do so, click the **Fetch** button, 
paste the recipe's URL and click *OK*. 

The application will display the recipe if the request is successful. Otherwise, you will be 
presented with a message asking you to either go back to the previous page or request the Recipya 
developers to support the website.

To view all supported websites, please click on the <ins>supported</ins> word.
You can scrape a website not in the supported list, but recipe extraction may fail. If it does, you can 
request support for the website by clicking the button that appears.

![img alt](/img/features/add-recipe-website.png)

### Import

You can import recipes in the `.json` that adhere to the [Recipe schema](https://schema.org/Recipe) standard. 

![img alt](/img/features/add-recipe-import.png)

You can upload either a single `.json` file or a zip file containing multiple recipes.
The recipes in a zip file may be organized by folder. Each folder may contain the `.json` recipe file and an image 
file. All other files in a folder will be ignored during processing. Here is an 
[example](https://sea.musicavis.ca/f/683b9b9a7cc84e1bac0c/?dl=1) of how such zip may look like.
