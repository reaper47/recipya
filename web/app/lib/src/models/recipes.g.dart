// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'recipes.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

RecipesModel _$RecipesModelFromJson(Map<String, dynamic> json) {
  return RecipesModel(
    recipes: (json['recipes'] as List<dynamic>)
        .map((e) => RecipeModel.fromJson(e as Map<String, dynamic>))
        .toList(),
  );
}

Map<String, dynamic> _$RecipesModelToJson(RecipesModel instance) =>
    <String, dynamic>{
      'recipes': instance.recipes,
    };
