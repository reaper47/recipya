extension Case on String? {
  String toTitleCase() {
    return this!
        .replaceAllMapped(
            RegExp(
              r'[A-Z]{2,}(?=[A-Z][a-z]+[0-9]*|\b)|[A-Z]?[a-z]+[0-9]*|[A-Z]|[0-9]+',
            ),
            (Match m) =>
                "${m[0]![0].toUpperCase()}${m[0]!.substring(1).toLowerCase()}")
        .replaceAll(RegExp(r'(_|-)+'), ' ');
  }

  String separateDelimeter(String delim) {
    var r = '([A-Za-z]*$delim)+';
    return this!.replaceAllMapped(RegExp(r), (Match m) => '${m[0]} ');
  }
}
