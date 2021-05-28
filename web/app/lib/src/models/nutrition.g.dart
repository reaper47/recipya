// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'nutrition.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

NutritionModel _$NutritionModelFromJson(Map<String, dynamic> json) {
  return NutritionModel(
    calories: json['calories'] as String?,
    carbohydrateContent: json['carbohydrateContent'] as String?,
    fatContent: json['fatContent'] as String?,
    saturatedFatContent: json['saturatedFatContent'] as String?,
    cholesterolContent: json['cholesterolContent'] as String?,
    proteinContent: json['proteinContent'] as String?,
    sodiumContent: json['sodiumContent'] as String?,
    fiberContent: json['fiberContent'] as String?,
    sugarContent: json['sugarContent'] as String?,
  );
}

Map<String, dynamic> _$NutritionModelToJson(NutritionModel instance) =>
    <String, dynamic>{
      'calories': instance.calories,
      'carbohydrateContent': instance.carbohydrateContent,
      'fatContent': instance.fatContent,
      'saturatedFatContent': instance.saturatedFatContent,
      'cholesterolContent': instance.cholesterolContent,
      'proteinContent': instance.proteinContent,
      'sodiumContent': instance.sodiumContent,
      'fiberContent': instance.fiberContent,
      'sugarContent': instance.sugarContent,
    };
