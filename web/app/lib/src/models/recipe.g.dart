// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'recipe.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

RecipeModel _$RecipeModelFromJson(Map<String, dynamic> json) {
  return RecipeModel(
    id: json['id'] as int?,
    name: json['name'] as String?,
    description: json['description'] as String?,
    url: json['url'] as String?,
    image: json['image'] as String?,
    prepTime: json['prepTime'] as String?,
    cookTime: json['cookTime'] as String?,
    totalTime: json['totalTime'] as String?,
    keywords: json['keywords'] as String?,
    recipeYield: json['recipeYield'] as int?,
    tool: (json['tool'] as List<dynamic>).map((e) => e as String?).toList(),
    recipeIngredients: (json['recipeIngredients'] as List<dynamic>)
        .map((e) => e as String?)
        .toList(),
    recipeInstructions: (json['recipeInstructions'] as List<dynamic>)
        .map((e) => e as String?)
        .toList(),
    nutrition: json['nutrition'] == null
        ? null
        : NutritionModel.fromJson(json['nutrition'] as Map<String, dynamic>),
    dateModified: json['dateModified'] as String?,
    dateCreated: json['dateCreated'] as String?,
  );
}

Map<String, dynamic> _$RecipeModelToJson(RecipeModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'description': instance.description,
      'url': instance.url,
      'image': instance.image,
      'prepTime': instance.prepTime,
      'cookTime': instance.cookTime,
      'totalTime': instance.totalTime,
      'keywords': instance.keywords,
      'recipeYield': instance.recipeYield,
      'tool': instance.tool,
      'recipeIngredients': instance.recipeIngredients,
      'recipeInstructions': instance.recipeInstructions,
      'nutrition': instance.nutrition,
      'dateModified': instance.dateModified,
      'dateCreated': instance.dateCreated,
    };
