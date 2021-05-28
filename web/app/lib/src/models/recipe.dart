import 'package:equatable/equatable.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import 'nutrition.dart';

part 'recipe.g.dart';

@JsonSerializable()
class RecipeModel extends Equatable {
  final int? id;
  final String? name;
  final String? description;
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

  factory RecipeModel.fromJson(Map<String, dynamic> json) =>
      _$RecipeModelFromJson(json);

  Map<String, dynamic> toJson() => _$RecipeModelToJson(this);

  @override
  List<Object?> get props => [
        id,
        name,
        description,
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
}
