part of 'recipes_provider.dart';

class RecipesNotifier extends StateNotifier<RecipesState> {
  final IRecipesRepository _repository;

  RecipesNotifier({
    required IRecipesRepository repository,
  })  : _repository = repository,
        super(const RecipesState.initial());

  Future<void> getSearchRecipes(String? ingredients, int? mode, int? n) async {
    state = const RecipesState.loading();

    try {
      final recipes = await _repository.getSearchRecipes(ingredients, mode, n);
      state = RecipesState.data(recipes: recipes);
    } catch (err) {
      state = RecipesState.error('Error fetching recipes: ${err.toString()}');
    }
  }

  void resetState() {
    state = RecipesState.initial();
  }
}
