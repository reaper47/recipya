import 'package:flutter/material.dart';
import 'package:recipe_hunter/src/views/search_args.dart';

const MODE1 = "Maximize items taken from the fridge";
const MODE2 = "Minimize number of items to buy";

class SearchPage extends StatefulWidget {
  static const routeName = '/';

  @override
  _SearchPageState createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  final _formKey = GlobalKey<FormState>();

  String _searchMode = MODE2;
  int _numRecipes = 10;

  List<TextEditingController> _ingredientControllers = [];
  List<Widget> _ingredientFields = [];
  List<FocusNode> _nodes = [];

  @override
  void initState() {
    super.initState();
    _ingredientFields.add(_createIngredientField());
  }

  @override
  void dispose() {
    _ingredientControllers.forEach((element) => element.dispose());
    _nodes.forEach((element) => element.dispose());
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Search Recipes'),
        actions: [
          TextButton(
            onPressed: () {
              if (_formKey.currentState!.validate()) {
                var mode = _searchMode == MODE1 ? 1 : 2;

                Navigator.pushNamed(
                  context,
                  '/results',
                  arguments: SearchArgs(_getIngredients(), _numRecipes, mode),
                );
              }
            },
            child: Text('Search'),
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: Padding(
          padding: const EdgeInsets.all(32.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Row(
                children: [
                  Expanded(
                    child: ListTile(
                      dense: true,
                      title: Text('Search Mode:'),
                      subtitle: DropdownButtonFormField<String>(
                        items: <String>[MODE1, MODE2]
                            .map((String value) => DropdownMenuItem<String>(
                                  value: value,
                                  child: Text(value),
                                ))
                            .toList(),
                        onChanged: (value) => _searchMode = value as String,
                        value: _searchMode,
                      ),
                    ),
                  ),
                  Expanded(
                    child: ListTile(
                      dense: true,
                      title: Text(
                        'Number of Recipes:',
                      ),
                      subtitle: DropdownButtonFormField<String>(
                        items: <int>[5, 10, 15, 20, 25, 30]
                            .map((int value) => DropdownMenuItem<String>(
                                  value: value.toString(),
                                  child: Text(value.toString()),
                                ))
                            .toList(),
                        onChanged: (value) => _numRecipes = int.parse(value!),
                        value: _numRecipes.toString(),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: 30,
              ),
              Expanded(
                child: ListView.builder(
                  itemCount: _ingredientFields.length,
                  itemBuilder: (BuildContext context, int index) => ListTile(
                    title: _ingredientFields[index],
                    trailing: GestureDetector(
                      onTapDown: (details) => _deleteEntry(index),
                      child: MouseRegion(
                        cursor: SystemMouseCursors.click,
                        child: Icon(Icons.delete),
                      ),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  TextFormField _createIngredientField() {
    var controller = TextEditingController();
    _ingredientControllers.add(controller);
    _nodes.add(FocusNode());
    return TextFormField(
      controller: controller,
      decoration: InputDecoration(hintText: 'Add an ingredient'),
      focusNode: _nodes[_nodes.length - 1],
      onEditingComplete: () {
        _onEditingCompleteCallback();
        _nodes.forEach((element) => element.unfocus());
        _nodes[_nodes.length - 1].requestFocus();
      },
      validator: (value) {
        if (value!.isEmpty) {
          return 'Please add an ingredient';
        }
      },
    );
  }

  void _onEditingCompleteCallback() {
    if (_ingredientControllers.where((el) => el.text == '').length == 0) {
      setState(() => _ingredientFields.add(_createIngredientField()));
    }
  }

  void _deleteEntry(int index) {
    if (_ingredientFields.length > 1) {
      setState(() {
        _ingredientControllers.removeAt(index);
        _ingredientFields.removeAt(index);
      });
    }
  }

  String _getIngredients() {
    return _ingredientControllers.map((e) => e.text).join(',');
  }
}

