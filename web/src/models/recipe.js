export default class Recipe {
  constructor(data) {
    this.category = data.recipeCategory;
    this.cookTime = data.cookTime;
    this.dateCreated = data.dateCreated;
    this.dateModified = data.dateModified;
    this.description = data.description;
    this.id = data.id;
    this.image = data.image;
    this.ingredients = data.recipeIngredient;
    this.instructions = data.recipeInstructions;
    this.keywords = data.keywords;
    this.name = data.name;
    this.nutrition = data.nutrition;
    this.prepTime = data.prepTime;
    this.tools = data.tool;
    this.totalTime = data.totalTime;
    this.url = data.url;
    this.yield = data.recipeYield;
  }
}
