import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:recipe_hunter/src/logic/states/recipes.dart';
import 'package:recipe_hunter/src/repository/recipes.dart';

part 'recipes_state_notifier.dart';

final recipesNotifierProvider =
    StateNotifierProvider<RecipesNotifier, RecipesState>(
  (ref) => RecipesNotifier(
    repository: ref.watch(_recipesRepositoryProvider),
  ),
);

final _recipesRepositoryProvider = Provider<IRecipesRepository>(
  (ref) => RecipesRepository(),
);
