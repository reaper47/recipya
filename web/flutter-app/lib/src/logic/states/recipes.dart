import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:recipe_hunter/src/models/recipes.dart';

part 'recipes.freezed.dart';

extension RecipesGetters on RecipesState {
  bool get isLoading => this is _RecipesStateLoading;
}

@freezed
class RecipesState with _$RecipesState {
  const factory RecipesState.initial() = _RecipesStateInitial;

  const factory RecipesState.loading() = _RecipesStateLoading;

  const factory RecipesState.data({required RecipesModel recipes}) =
      _RecipesStateData;

  const factory RecipesState.error([String? error]) = _RecipesStateError;
}
