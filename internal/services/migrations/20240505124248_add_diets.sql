-- +goose Up
CREATE TABLE diets
(
    id          INTEGER PRIMARY KEY,
    category    TEXT NOT NULL,
    name        TEXT NOT NULL UNIQUE,
    schema_link INTEGER REFERENCES diet_schema (id) ON DELETE CASCADE
);

CREATE TABLE diet_category (
    id INTEGER PRIMARY KEY ,
    name TEXT NOT NULL
);

CREATE TABLE diet_schema (
  id INTEGER PRIMARY KEY ,
  name TEXT NOT NULL,
  schema TEXT NOT  NULL,
  UNIQUE (name, schema)
);

INSERT INTO diet_category (name)
VALUES ("Belief-based"),
       ("Crash"),
       ("Detox"),
       ("Fad - Food-specific"),
       ("Fad - Low-carbohydrate / high-fat"),
       ("Fad - High-carbohydrate / low-fat "),
       ("Fad - Liquid"),
       ("Fad - Fasting"),
       ("Fad - Detoxifying"),
       ("Fad - Other"),
       ("Medical reasons"),
       ("Other"),
       ("Vegetarian"),
       ("Vegetarian - Semi-vegetarian"),
       ("Weight contol - Low-calorie"),
       ("Weight contol - Very low calorie"),
       ("Weight contol - Low-carbohydrate"),
       ("Weight contol - Low-fat");

INSERT INTO diet_schema (name, schema)
VALUES ("Diabetic", "https://schema.org/DiabeticDiet"),
       ("Gluten-Free", "https://schema.org/GlutenFreeDiet"),
       ("Halal", "https://schema.org/HalalDiet"),
       ("Hindu", "https://schema.org/HinduDiet"),
       ("Kosher", "https://schema.org/KosherDiet"),
       ("Low Calorie", "https://schema.org/LowCalorieDiet"),
       ("Low Fat", "https://schema.org/LowFatDiet"),
       ("Low Lactose", "https://schema.org/LowLactoseDiet"),
       ("Low Salt", "https://schema.org/LowSaltDiet"),
       ("Vegan", "https://schema.org/VeganDiet"),
       ("Vegetarian", "https://schema.org/VegetarianDiet");

INSERT INTO diets (category, name, schema_link)
VALUES ((SELECT id FROM diet_category WHERE name = "Belief-based"), "Buddhist", NULL),
       ((SELECT id FROM diet_category WHERE name = "Belief-based"), "Halal", (SELECT id FROM diet_schema WHERE name = "Halal")),
       ((SELECT id FROM diet_category WHERE name = "Belief-based"), "Hindu", (SELECT id FROM diet_schema WHERE name = "Hindu")),
       ((SELECT id FROM diet_category WHERE name = "Belief-based"), "Jain", NULL),
       ((SELECT id FROM diet_category WHERE name = "Belief-based"), "I-tal", NULL),
       ((SELECT id FROM diet_category WHERE name = "Belief-based"), "Kosher", (SELECT id FROM diet_schema WHERE name = "Kosher")),
       ((SELECT id FROM diet_category WHERE name = "Belief-based"), "Seventh-day Adventist", NULL),
       ((SELECT id FROM diet_category WHERE name = "Belief-based"), "Word of Wisdom", NULL),
       ((SELECT id FROM diet_category WHERE name = "Crash"), "General", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Crash"), "Beverly Hills", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Crash"), "Cabbage Soup", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Crash"), "Grapefruit", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Crash"), "Monotrophic", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Crash"), "Subway", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Detox"), "General", NULL),
       ((SELECT id FROM diet_category WHERE name = "Detox"), "Juice Fasting", NULL),
       ((SELECT id FROM diet_category WHERE name = "Detox"), "Master Cleanse", NULL),
       ((SELECT id FROM diet_category WHERE name = "Medical reasons"), "General", NULL),
       ((SELECT id FROM diet_category WHERE name = "Medical reasons"), "DASH", NULL),
       ((SELECT id FROM diet_category WHERE name = "Medical reasons"), "Diabetic", (SELECT id FROM diet_schema WHERE name = "Diabetic")),
       ((SELECT id FROM diet_category WHERE name = "Medical reasons"), "Elemental", NULL),
       ((SELECT id FROM diet_category WHERE name = "Medical reasons"), "Elimination", NULL),
       ((SELECT id FROM diet_category WHERE name = "Medical reasons"), "Gluten-free", (SELECT id FROM diet_schema WHERE name = "Gluten-Free")),
       ((SELECT id FROM diet_category WHERE name = "Medical reasons"), "Gluten-free, casein-free", (SELECT id FROM diet_schema WHERE name = "Gluten-Free")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "General", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "5:2", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "Intermitten fasting", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "Body for Life", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "Cookie", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "The Hacker", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "Nutrisystem", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-calorie"), "Weight Watchers", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Very low calorie"), "General", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Very low calorie"), "Inedia", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Very low calorie"), "KE", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Very low calorie"), "The Last Chance", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Very low calorie"), "Tongue Patch", (SELECT id FROM diet_schema WHERE name = "Low Calorie")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-carbohydrate"), "General", NULL),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-carbohydrate"), "Atkins", NULL),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-carbohydrate"), "Dukan", NULL),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-carbohydrate"), "Kimkins", NULL),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-carbohydrate"), "South Beach", NULL),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-carbohydrate"), "Stillman", NULL),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-fat"), "General", (SELECT id FROM diet_schema WHERE name = "Low Fat")),
       ((SELECT id FROM diet_category WHERE name = "Weight contol - Low-fat"), "McDougall's starch", (SELECT id FROM diet_schema WHERE name = "Low Fat")),

-- +goose Down

