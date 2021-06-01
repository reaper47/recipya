import 'package:equatable/equatable.dart';

class NutritionModel extends Equatable {
  final String? calories;
  final String? carbohydrateContent;
  final String? fatContent;
  final String? saturatedFatContent;
  final String? cholesterolContent;
  final String? proteinContent;
  final String? sodiumContent;
  final String? fiberContent;
  final String? sugarContent;

  NutritionModel({
    this.calories,
    this.carbohydrateContent,
    this.fatContent,
    this.saturatedFatContent,
    this.cholesterolContent,
    this.proteinContent,
    this.sodiumContent,
    this.fiberContent,
    this.sugarContent,
  });

  @override
  List<Object?> get props => [
        calories,
        carbohydrateContent,
        fatContent,
        saturatedFatContent,
        cholesterolContent,
        proteinContent,
        sodiumContent,
        fiberContent,
        sugarContent,
      ];

  NutritionModel.fromJson(Map<String, dynamic> json)
      : calories = json['calories'],
        carbohydrateContent = json['carbohydrate'],
        fatContent = json['fat'],
        saturatedFatContent = json['saturatedFat'],
        cholesterolContent = json['cholesterol'],
        proteinContent = json['protein'],
        sodiumContent = json['sodium'],
        fiberContent = json['fiber'],
        sugarContent = json['sugar'];
}
