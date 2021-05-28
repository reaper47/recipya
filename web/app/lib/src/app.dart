import 'package:flutter/material.dart';

import 'views/search.dart';

class RecipeHunterApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Recipe Hunter',
      theme: ThemeData.dark().copyWith(
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: SearchPage(),
    );
  }
}
