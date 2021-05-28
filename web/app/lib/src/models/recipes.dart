import 'package:equatable/equatable.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import 'recipe.dart';

part 'recipes.g.dart';

@JsonSerializable()
class RecipesModel extends Equatable {
  final List<RecipeModel> recipes;

  RecipesModel({required this.recipes});

  @override
  List<Object?> get props => [
        recipes,
      ];

  factory RecipesModel.fromJson(Map<String, dynamic> json) =>
      _$RecipesModelFromJson(json);

  Map<String, dynamic> toJson() => _$RecipesModelToJson(this);
}
