import 'dart:convert';

import 'package:http/http.dart' as http;

import 'package:recipe_hunter/src/models/recipes.dart';

abstract class IRecipesRepository {
  Future<RecipesModel> getSearchRecipes(String? ingredients, int? mode, int? n);
}

class RecipesRepository implements IRecipesRepository {
  String baseUrl = Uri.base.origin + "/api/v1";

  @override
  Future<RecipesModel> getSearchRecipes(
      String? ingredients, int? mode, int? n) async {
    try {
      var url =
          Uri.parse('$baseUrl/search?ingredients=$ingredients&mode=$mode&n=$n');
      final response = await http.get(url);
      if (response.statusCode == 200) {
        return RecipesModel.fromJson(jsonDecode(response.body));
      } else {
        throw Exception();
      }
    } catch (e) {
      throw Exception();
    }
  }
}
