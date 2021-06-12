import 'package:equatable/equatable.dart';
import 'recipe.dart';

class RecipesModel extends Equatable {
  final List<RecipeModel> recipes;

  RecipesModel({required this.recipes});

  @override
  List<Object?> get props => [
        recipes,
      ];

  RecipesModel.fromJson(Map<String, dynamic> json)
      : recipes = [
          for (var recipe in json['recipes']) RecipeModel.fromJson(recipe)
        ];
}
