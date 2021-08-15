import argparse
from datetime import datetime
import json
import sys
from urllib.parse import urlparse

from recipe_scrapers import scrape_me
import scrape_schema_recipe


def main() -> None:
    parser = init_argparse()
    url = parser.parse_args().url
    now = str(datetime.now())

    try:
        scraper = scrape_me(url)
        schema = scrape_schema_recipe.loads(str(scraper.soup))

        name = scraper.title()
        image = scraper.image()
        yields = get_yield(scraper)
        recipeIngredient = scraper.ingredients()
        recipeInstructions = scraper.instructions().split("\n")
        nutrition = scraper.nutrients()

    except:
        schema = scrape_schema_recipe.scrape_url(url)
        if len(schema) == 0:
            sys.exit(1)

        name = get_field_from_schema("name", "", schema)
        image = get_field_from_schema("image", "", schema)[0]
        yields = get_field_from_schema("recipeYield", "", schema)[0]
        recipeIngredient = get_field_from_schema("recipeIngredient", [], schema)
        recipeInstructions = get_field_from_schema("recipeInstructions", [], schema)
        recipeInstructions = get_instructions_from_schema(recipeInstructions)
        nutrition = get_field_from_schema("nutrients", {}, schema)

    recipe = json.dumps(
        {
            "name": name,
            "description": get_field_from_schema("description", "", schema),
            "url": url,
            "image": image,
            "prepTime": get_field_from_schema("prepTime", "PT0H", schema),
            "cookTime": get_field_from_schema("cookTime", "PT0H", schema),
            "totalTime": get_field_from_schema("totalTime", "PT0H", schema),
            "recipeCategory": get_field_from_schema("recipeCategory", [""], schema)[0],
            "keywords": get_field_from_schema("keywords", "", schema),
            "recipeYield": yields,
            "tool": get_field_from_schema("tool", [], schema),
            "recipeIngredient": recipeIngredient,
            "recipeInstructions": recipeInstructions,
            "nutrition": nutrition,
            "dateModified": now,
            "dateCreated": now,
        }
    )
    print(recipe)


def init_argparse() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        usage="%(prog)s [URL]",
        description="Get recipe data from a URL and return it in a JSON object.",
    )
    parser.add_argument("url", type=url, help="the URL of the recipe to import")
    return parser


def url(str):
    if not is_url(str):
        raise ValueError()
    return str


def is_url(url):
    try:
        result = urlparse(url)
        return all([result.scheme, result.netloc])
    except ValueError:
        return False


def get_field_from_schema(field, default, schema):
    if len(schema) == 0:
        return default

    recipe = schema[0]
    if field not in recipe:
        return default

    return recipe[field]


def get_yield(scraper):
    yields = scraper.yields()
    if type(yields) is str:
        for part in yields.split(" "):
            try:
                return int(part)
            except ValueError:
                continue
    return 0


def get_instructions_from_schema(instructions):
    steps = []
    for instruction in instructions:
        name = f"<{instruction['name']}>"
        for step in instruction["itemListElement"]:
            steps.append(f"{name}{step['text']}")
    return steps


if __name__ == "__main__":
    main()
