import 'package:equatable/equatable.dart';

import 'nutrition.dart';

class RecipeModel extends Equatable {
  final int? id;
  final String? name;
  final String? description;
  final String? category;
  final String? url;
  final String? image;
  final String? prepTime;
  final String? cookTime;
  final String? totalTime;
  final String? keywords;
  final int? recipeYield;
  final List<String?> tool;
  final List<String?> recipeIngredients;
  final List<String?> recipeInstructions;
  final NutritionModel? nutrition;
  final String? dateModified;
  final String? dateCreated;

  RecipeModel({
    this.id,
    this.name,
    this.description,
    this.category,
    this.url,
    this.image,
    this.prepTime,
    this.cookTime,
    this.totalTime,
    this.keywords,
    this.recipeYield,
    required this.tool,
    required this.recipeIngredients,
    required this.recipeInstructions,
    this.nutrition,
    this.dateModified,
    this.dateCreated,
  });

  @override
  List<Object?> get props => [
        id,
        name,
        description,
        category,
        url,
        image,
        prepTime,
        cookTime,
        totalTime,
        keywords,
        recipeYield,
        tool,
        recipeIngredients,
        recipeInstructions,
        nutrition,
        dateModified,
        dateCreated,
      ];

  RecipeModel.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        name = json['name'],
        description = json['description'],
        category = json['recipeCategory'],
        url = json['url'],
        image = json['image'],
        prepTime = json['prepTime'],
        cookTime = json['cookTime'],
        totalTime = json['totalTime'],
        keywords = json['keywords'],
        recipeYield = json['recipeYield'],
        tool = json['tool'] ?? [],
        recipeIngredients = (json['recipeIngredient'] as List).cast<String>(),
        recipeInstructions =
            (json['recipeInstructions'] as List).cast<String>(),
        nutrition = NutritionModel.fromJson(json['nutrition']),
        dateModified = json['dateModified'],
        dateCreated = json['dateCreated'];

  String totalTimeToString() {
    if (totalTime == "" || totalTime == "PT0H0M") {
      return "0m";
    }

    var parts = totalTime!.split('PT')[1].split('H');
    var hour = parts[0];
    var min = parts[1].split('M')[0];
    return '${hour}h${min}m';
  }
}
