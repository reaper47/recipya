import 'package:flutter/material.dart';
import 'package:recipe_hunter/src/views/search_results.dart';

import 'views/search.dart';

class RecipeHunterApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Recipe Hunter',
      theme: ThemeData.dark().copyWith(
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      initialRoute: SearchPage.routeName,
      routes: {
        SearchPage.routeName: (context) => SearchPage(),
        SearchResultsPage.routeName: (context) => SearchResultsPage(),
      },
    );
  }
}
