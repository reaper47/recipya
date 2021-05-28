import 'package:equatable/equatable.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'nutrition.g.dart';

@JsonSerializable()
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

  factory NutritionModel.fromJson(Map<String, dynamic> json) =>
      _$NutritionModelFromJson(json);

  Map<String, dynamic> toJson() => _$NutritionModelToJson(this);
}
