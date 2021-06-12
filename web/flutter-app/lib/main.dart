import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:recipe_hunter/src/app.dart';

void main() {
  runApp(
    ProviderScope(
      child: RecipeHunterApp(),
    ),
  );
}
