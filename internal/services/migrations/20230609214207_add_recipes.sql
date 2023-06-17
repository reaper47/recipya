-- +goose Up
--
-- Create normal tables
--
CREATE TABLE recipes
(
    id             INTEGER PRIMARY KEY,
    name           TEXT      NOT NULL,
    description    TEXT,
    image          TEXT               DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' ||
                                               substr(lower(hex(randomblob(2))), 2) || '-a' ||
                                               substr(lower(hex(randomblob(2))), 2) || '-%' ||
                                               substr(lower(hex(randomblob(6))), 2)),
    yield          INTEGER,
    url            TEXT,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- https://recipeland.com/recipes/categories/browse
CREATE TABLE categories
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE cuisines
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE ingredients
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE instructions
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE tools
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE keywords
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE nutrition
(
    id                  INTEGER PRIMARY KEY,
    recipe_id           INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    calories            TEXT,
    total_carbohydrates TEXT,
    sugars              TEXT,
    protein             TEXT,
    total_fat           TEXT,
    saturated_fat       TEXT,
    unsaturated_fat     TEXT,
    cholesterol         TEXT,
    sodium              TEXT,
    fiber               TEXT
);

CREATE TABLE times
(
    id            INTEGER PRIMARY KEY,
    prep_seconds  INTEGER DEFAULT 0,
    cook_seconds  INTEGER DEFAULT 0,
    total_seconds INTEGER DEFAULT NULL,
    UNIQUE (prep_seconds, cook_seconds)
);

CREATE TABLE counts
(
    id      INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    recipes INTEGER DEFAULT 0
);

CREATE TABLE fruits_veggies
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE blacklist_units
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

--
-- Create association Tables
--
CREATE TABLE user_recipe
(
    id        INTEGER PRIMARY KEY,
    user_id   INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (user_id, recipe_id)
);

CREATE TABLE user_category
(
    id          INTEGER PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    UNIQUE (user_id, category_id)
);

CREATE TABLE category_recipe
(
    id          INTEGER PRIMARY KEY,
    category_id INTEGER DEFAULT 1 REFERENCES categories (id) ON DELETE SET DEFAULT,
    recipe_id   INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (category_id, recipe_id)
);

CREATE TABLE cuisine_recipe
(
    id         INTEGER PRIMARY KEY,
    cuisine_id INTEGER DEFAULT 1 REFERENCES cuisines (id) ON DELETE SET DEFAULT,
    recipe_id  INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (cuisine_id, recipe_id)
);

CREATE TABLE ingredient_recipe
(
    id            INTEGER PRIMARY KEY,
    ingredient_id INTEGER NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    recipe_id     INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, ingredient_id)
);

CREATE TABLE instruction_recipe
(
    id             INTEGER PRIMARY KEY,
    instruction_id INTEGER NOT NULL REFERENCES instructions (id) ON DELETE CASCADE,
    recipe_id      INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, instruction_id)
);

CREATE TABLE keyword_recipe
(
    id         INTEGER PRIMARY KEY,
    keyword_id INTEGER NOT NULL REFERENCES keywords (id) ON DELETE CASCADE,
    recipe_id  INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, keyword_id)
);

CREATE TABLE tool_recipe
(
    id        INTEGER PRIMARY KEY,
    tool_id   INTEGER NOT NULL REFERENCES tools (id) ON DELETE CASCADE,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, tool_id)
);


CREATE TABLE time_recipe
(
    id        INTEGER PRIMARY KEY,
    time_id   INTEGER NOT NULL REFERENCES times (id) ON DELETE SET NULL,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, time_id)
);

--
-- FUNCTIONS
--

-- +goose StatementBegin
CREATE TRIGGER users_insert_trigger
    AFTER INSERT
    ON users
    FOR EACH ROW
BEGIN
    INSERT INTO counts (user_id)
    VALUES (NEW.id);

    INSERT INTO user_category (user_id, category_id)
    SELECT NEW.id, id
    FROM categories
    WHERE name IN (
                   'uncategorized', 'appetizers', 'bread', 'breakfasts', 'condiments',
                   'dessert', 'lunch', 'main dish', 'salad', 'side dish',
                   'snacks', 'soups', 'stews'
        );
END;

CREATE TRIGGER times_calc_total_time_ai
    AFTER INSERT
    ON times
    FOR EACH ROW
    WHEN (NEW.total_seconds IS NULL)
BEGIN
    UPDATE times
    SET total_seconds = prep_seconds + times.cook_seconds
    WHERE id = NEW.id;
END;

CREATE TRIGGER times_calc_total_time_bu
    AFTER UPDATE
    ON times
    FOR EACH ROW
BEGIN
    UPDATE times
    SET total_seconds = prep_seconds + times.cook_seconds
    WHERE id = NEW.id;
END;

CREATE TRIGGER recipes_update_updated_at
    AFTER UPDATE
    ON recipes
    FOR EACH ROW
BEGIN
    UPDATE recipes
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

CREATE TRIGGER user_recipe_insert_inc_count
    AFTER INSERT
    ON user_recipe
    FOR EACH ROW
BEGIN
    UPDATE counts
    SET recipes = recipes + 1
    WHERE id = NEW.user_id;
END;

CREATE TRIGGER user_recipe_delete_dec_count
    AFTER DELETE
    ON user_recipe
    FOR EACH ROW
BEGIN
    UPDATE counts
    SET recipes = recipes - 1
    WHERE id = OLD.user_id;
END;
-- +goose StatementEnd

--
-- INSERTS
--
INSERT INTO categories (name)
VALUES ('uncategorized'),
       ('appetizers'),
       ('bread'),
       ('breakfasts'),
       ('condiments'),
       ('dessert'),
       ('lunch'),
       ('main dish'),
       ('salad'),
       ('side dish'),
       ('snacks'),
       ('soups'),
       ('stews');

INSERT INTO fruits_veggies (name)
VALUES ('abiu'),
       ('acerola'),
       ('ackee'),
       ('acorn squash'),
       ('adzuki'),
       ('african cucumber'),
       ('alfalfa sprouts'),
       ('anise'),
       ('apple'),
       ('apricot'),
       ('artichoke'),
       ('arugula'),
       ('asparagus'),
       ('aubergine'),
       ('avocado'),
       ('azuki beans'),
       ('açaí'),
       ('banana'),
       ('banana squash'),
       ('basil'),
       ('bean sprouts'),
       ('beet'),
       ('beet greens'),
       ('beetroot'),
       ('bell pepper'),
       ('bilberry'),
       ('bitter melon'),
       ('black beans'),
       ('black sapote'),
       ('black-eyed peas'),
       ('blackberry'),
       ('blackcurrant'),
       ('blood orange'),
       ('blueberry'),
       ('bok choy'),
       ('borlotti bean'),
       ('boysenberry'),
       ('breadfruit'),
       ('broad beans'),
       ('broccoflower'),
       ('broccoli'),
       ('brussels sprouts'),
       ('buddhas hand'),
       ('butter bean'),
       ('butternut squash'),
       ('cabbage'),
       ('cactus pear'),
       ('calabrese'),
       ('canistel'),
       ('cantaloupe'),
       ('caraway'),
       ('carrot'),
       ('cauliflower'),
       ('caviar lime'),
       ('cayenne pepper'),
       ('ceci beans'),
       ('celeriac'),
       ('celery'),
       ('cempedak'),
       ('chamomile'),
       ('chard'),
       ('cherimoya'),
       ('cherry'),
       ('chickpeas'),
       ('chico fruit'),
       ('chili pepper'),
       ('chives'),
       ('clementine'),
       ('cloudberry'),
       ('coco de mer'),
       ('coconut'),
       ('collard greens'),
       ('coriander'),
       ('corms'),
       ('courgette'),
       ('crab apple'),
       ('cranberry'),
       ('cucumber'),
       ('currant'),
       ('custard apple'),
       ('cymbopogon'),
       ('daikon'),
       ('damson'),
       ('date'),
       ('delicata'),
       ('dill'),
       ('dragonfruit'),
       ('dried plum'),
       ('durian'),
       ('eddoe'),
       ('egg fruit'),
       ('eggplant'),
       ('elderberry'),
       ('endive'),
       ('falsa'),
       ('feijoa'),
       ('fennel'),
       ('fiddleheads'),
       ('fig'),
       ('finger lime'),
       ('frisee'),
       ('galia melon'),
       ('garbanzos'),
       ('garlic'),
       ('gem squash'),
       ('ginger'),
       ('goji berry'),
       ('gooseberry'),
       ('grape'),
       ('grapefruit'),
       ('green beans'),
       ('green onion'),
       ('greens'),
       ('grewia asiatica'),
       ('guava'),
       ('habanero'),
       ('hala fruit'),
       ('hawthorn berry'),
       ('herbs'),
       ('honeyberry'),
       ('honeydew'),
       ('horned melon'),
       ('horseradish'),
       ('hubbard squash'),
       ('huckleberry'),
       ('jabuticaba'),
       ('jackfruit'),
       ('jalapeÃ±o'),
       ('jambul'),
       ('japanese plum'),
       ('jerusalem artichoke'),
       ('jicama'),
       ('jostaberry'),
       ('jujube'),
       ('juniper berry'),
       ('kaffir lime'),
       ('kale'),
       ('kidney beans'),
       ('kiwano'),
       ('kiwi'),
       ('kiwis'),
       ('kohlrabi'),
       ('konjac'),
       ('kumquat'),
       ('lavender'),
       ('leek'),
       ('legumes'),
       ('lemon'),
       ('lemongrass'),
       ('lentils'),
       ('lettuce'),
       ('lima beans'),
       ('lime'),
       ('loganberry'),
       ('longan'),
       ('loquat'),
       ('lulo'),
       ('lychee'),
       ('magellan barberry'),
       ('mamey apple'),
       ('mamey sapote'),
       ('mamin chino'),
       ('mandarine'),
       ('mangel-wurzel'),
       ('mangetout'),
       ('mango'),
       ('mangosteen'),
       ('marionberry'),
       ('marjoram'),
       ('marrow'),
       ('melon'),
       ('miracle fruit'),
       ('monstera deliciosa'),
       ('mouse melon'),
       ('mulberry'),
       ('mung beans'),
       ('mushrooms'),
       ('musk melon'),
       ('mustard greens'),
       ('nance'),
       ('navy beans'),
       ('nectarine'),
       ('nettles'),
       ('new zealand spinach'),
       ('okra'),
       ('onion'),
       ('orange'),
       ('oregano'),
       ('oyster plant'),
       ('papaya'),
       ('paprika'),
       ('parsley'),
       ('parsnip'),
       ('passionfruit'),
       ('pawpaw'),
       ('peach'),
       ('peanuts'),
       ('pear'),
       ('peas'),
       ('peppers'),
       ('persimmon'),
       ('phalsa'),
       ('pineapple'),
       ('pineberry'),
       ('pinto beans'),
       ('pitaya'),
       ('plantain'),
       ('plum'),
       ('plumcot'),
       ('pluot'),
       ('pomegranate'),
       ('pomelo'),
       ('potato'),
       ('prune'),
       ('purple mangosteen'),
       ('quince'),
       ('radicchio'),
       ('radish'),
       ('raisin'),
       ('rambutan'),
       ('raspberry'),
       ('red cabbage'),
       ('redcurrant'),
       ('rhubarb'),
       ('root vegetables'),
       ('rose apple'),
       ('rosemary'),
       ('runner beans'),
       ('rutabaga'),
       ('salak'),
       ('salal berry'),
       ('salmonberry'),
       ('salsify'),
       ('satsuma'),
       ('savoy cabbage'),
       ('scallion'),
       ('shallot'),
       ('shine muscat '),
       ('skirret'),
       ('sloe'),
       ('snap peas'),
       ('soursop'),
       ('soy beans'),
       ('spaghetti squash'),
       ('spinach'),
       ('split peas'),
       ('spring onion'),
       ('squash'),
       ('star apple'),
       ('star fruit'),
       ('strawberry'),
       ('surinam cherry'),
       ('sweet potato'),
       ('sweetcorn'),
       ('tabasco pepper'),
       ('tamarillo'),
       ('tamarind'),
       ('tangelo'),
       ('tangerine'),
       ('taro'),
       ('tat soi'),
       ('tayberry'),
       ('thyme'),
       ('tomato'),
       ('topinambur'),
       ('tubers'),
       ('turnip'),
       ('ugli fruit'),
       ('vitis vinifera'),
       ('wasabi'),
       ('water chestnut'),
       ('watercresswhite radish'),
       ('watermelon'),
       ('white currant'),
       ('white sapote'),
       ('yam'),
       ('yuzu'),
       ('zucchini');

INSERT INTO blacklist_units (name)
VALUES ('¼'),
       ('½'),
       ('¾'),
       ('bar'),
       ('c'),
       ('cc'),
       ('celsius'),
       ('centimeter'),
       ('centimetre'),
       ('cm'),
       ('cube'),
       ('cup'),
       ('deciliter'),
       ('decilitre'),
       ('dl'),
       ('fahrenheit'),
       ('fl'),
       ('fluid'),
       ('g'),
       ('gal'),
       ('gallon'),
       ('gill'),
       ('gram'),
       ('gramme'),
       ('imperial'),
       ('in'),
       ('inch'),
       ('kg'),
       ('kilogram'),
       ('kilogramme'),
       ('l'),
       ('lb'),
       ('liter'),
       ('litre'),
       ('m'),
       ('meter'),
       ('metre'),
       ('mg'),
       ('milligram'),
       ('milligramme'),
       ('milliliter'),
       ('millilitre'),
       ('millimeter'),
       ('millimetre'),
       ('ml'),
       ('mm'),
       ('optional'),
       ('ounce'),
       ('oz'),
       ('package'),
       ('pint'),
       ('pound'),
       ('quart'),
       ('tablespoon'),
       ('tbl'),
       ('tbs'),
       ('tbsp'),
       ('teaspoon'),
       ('tsp');

-- +goose Down
DROP TRIGGER user_recipe_delete_dec_count;
DROP TRIGGER user_recipe_insert_inc_count;
DROP TRIGGER recipes_update_updated_at;
DROP TRIGGER times_calc_total_time_bu;
DROP TRIGGER times_calc_total_time_ai;
DROP TRIGGER users_insert_trigger;

DROP TABLE time_recipe;
DROP TABLE keyword_recipe;
DROP TABLE tool_recipe;
DROP TABLE instruction_recipe;
DROP TABLE ingredient_recipe;
DROP TABLE category_recipe;
DROP TABLE user_recipe;
DROP TABLE blacklist_units;
DROP TABLE fruits_veggies;
DROP TABLE user_category;
DROP TABLE counts;
DROP TABLE times;
DROP TABLE nutrition;
DROP TABLE keywords;
DROP TABLE tools;
DROP TABLE instructions;
DROP TABLE ingredients;
DROP TABLE categories;
DROP TABLE recipes;
