import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:recipe_hunter/src/logic/providers/recipes_provider.dart';
import 'package:recipe_hunter/src/models/recipe.dart';
import 'package:recipe_hunter/src/views/search_args.dart';
import 'package:recipe_hunter/src/util/strings_ext.dart';

class SearchResultsPage extends StatefulHookWidget {
  static const routeName = '/results';

  const SearchResultsPage({Key? key}) : super(key: key);

  @override
  _SearchResultsPageState createState() => _SearchResultsPageState();
}

class _SearchResultsPageState extends State<SearchResultsPage> {
  @override
  void initState() {
    Future.delayed(Duration.zero, () {
      final args = ModalRoute.of(context)!.settings.arguments as SearchArgs;
      context
          .read(recipesNotifierProvider.notifier)
          .getSearchRecipes(args.ingredients, args.mode, args.limit);
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          onPressed: () {
            Navigator.pop(context);
          },
          icon: Icon(Icons.arrow_back),
        ),
        title: Text('Search Results'),
      ),
      body: Padding(
        padding: EdgeInsets.symmetric(horizontal: 20.0),
        child: _RecipesConsumer(),
      ),
    );
  }
}

class _RecipesConsumer extends ConsumerWidget {
  @override
  Widget build(BuildContext context, ScopedReader watch) {
    final state = watch(recipesNotifierProvider);

    return state.when(
      initial: () => Text('Something happened :/ Please try again.'),
      loading: () => Center(
        child: CircularProgressIndicator(),
      ),
      data: (data) {
        var recipes = data.recipes;
        if (recipes.isEmpty) {
          return Center(
            child: Text('No recipes found for this query.'),
          );
        }
        return GridView.count(
          physics: BouncingScrollPhysics(),
          padding: const EdgeInsets.all(32),
          crossAxisCount: 3,
          children: [for (var recipe in recipes) _createRecipeCard(recipe)],
        );
      },
      error: (error) => Text('$error'),
    );
  }

  Widget _createRecipeCard(RecipeModel recipe) {
    return Card(
      child: InkWell(
        splashColor: Colors.blue.withAlpha(30),
        onTap: () {
          print('Card tapped.');
        },
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            children: [
              Text(
                '${recipe.name.toTitleCase()}',
                style: TextStyle(
                  fontSize: 16,
                  decoration: TextDecoration.underline,
                ),
              ),
              const SizedBox(
                height: 15,
              ),
              Container(
                height: 130,
                child: Text(
                  '${recipe.description}',
                  maxLines: 7,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(fontSize: 14, height: 1.25),
                ),
              ),
              Spacer(),
              ListTile(
                //isThreeLine: true,
                leading: Icon(Icons.food_bank),
                title: Text(
                  '${recipe.category.toTitleCase()}',
                  style: TextStyle(fontSize: 14),
                ),
                subtitle: Text(
                  '${recipe.keywords.separateDelimeter(',')}',
                  style: TextStyle(fontSize: 14),
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
