import 'package:dio/dio.dart';
import 'package:recipe_hunter/src/models/recipes.dart';

abstract class IRecipesRepository {
  Future<RecipesModel> getSearchRecipes(int? mode, int? n);
}

class RecipesRepository implements IRecipesRepository {
  final _dio = Dio();
  final baseUrl = "http://192.168.0.109:3000/api/v1/recipes";

  @override
  Future<RecipesModel> getSearchRecipes(int? mode, int? n) async {
    try {
      String url = '$baseUrl/search?ingredients=avocado&mode=$mode&n=$n';
      final result = await _dio.get(url);
      if (result.statusCode == 200) {
        return RecipesModel.fromJson(result.data);
      } else {
        throw Exception();
      }
    } catch (e) {
      throw Exception();
    }
  }
} // you are smart :)
