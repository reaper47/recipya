import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:recipe_hunter/src/logic/providers/recipes_provider.dart';
import 'package:recipe_hunter/src/logic/states/recipes.dart';

class SearchPage extends StatelessWidget {
  const SearchPage({Key? key}) : super(key: key);

  static const routeName = 'SearchPage';

  static Route route() {
    return MaterialPageRoute<void>(builder: (_) => const SearchPage());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Search Recipes'),
      ),
      body: Padding(
        padding: EdgeInsets.symmetric(horizontal: 20.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            _RecipesConsumer(),
            const SizedBox(height: 50),
            _ButtonConsumer(),
          ],
        ),
      ),
    );
  }
}

class _RecipesConsumer extends ConsumerWidget {
  @override
  Widget build(BuildContext context, ScopedReader watch) {
    final state = watch(recipesNotifierProvider);

    return state.when(
      initial: () => Text(
        'Press the button to start',
        textAlign: TextAlign.center,
      ),
      loading: () => Center(
        child: CircularProgressIndicator(),
      ),
      data: (recipes) => Text('$recipes'),
      error: (error) => Text('$error'),
    );
  }
}

class _ButtonConsumer extends ConsumerWidget {
  @override
  Widget build(BuildContext context, ScopedReader watch) {
    final state = watch(recipesNotifierProvider);

    return ElevatedButton(
      child: Text('Press me to get recipes'),
      onPressed: state != RecipesState.loading()
          ? () {
              context
                  .read(recipesNotifierProvider.notifier)
                  .getSearchRecipes(2, 2);
            }
          : null,
    );
  }
}
