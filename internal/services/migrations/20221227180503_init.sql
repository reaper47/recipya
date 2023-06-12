-- +goose Up

-- Tables
CREATE TABLE users
(
    id              INTEGER PRIMARY KEY,
    email           TEXT      NOT NULL UNIQUE CHECK ( LOWER(email) = email ),
    hashed_password TEXT      NOT NULL,
    is_confirmed    INTEGER            DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth_tokens
(
    id             INTEGER PRIMARY KEY,
    selector       CHAR(12),
    hash_validator CHAR(64),
    expires        INTEGER DEFAULT (unixepoch('now', '+1 month')),
    user_id        INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE recipes
(
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    url         TEXT,
    image       BLOB,
    yield       INTEGER,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- https://recipeland.com/recipes/categories/browse
CREATE TABLE categories
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
    calories            TEXT,
    total_carbohydrates TEXT,
    sugars              TEXT,
    protein             TEXT,
    total_fat           TEXT,
    saturated_fat       TEXT,
    cholesterol         TEXT,
    sodium              TEXT,
    fiber               TEXT,
    recipe_id           INTEGER REFERENCES recipes (id) ON DELETE CASCADE
);

CREATE TABLE times
(
    id    INTEGER PRIMARY KEY,
    prep  INTEGER,
    cook  INTEGER,
    total INTEGER,
    UNIQUE (prep, cook)
);

CREATE TABLE counts
(
    id          INTEGER PRIMARY KEY,
    num_recipes INTEGER DEFAULT 0,
    user_id     INTEGER REFERENCES users (id) ON DELETE CASCADE
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


-- Association Tables
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

CREATE TABLE recipe_category
(
    id          INTEGER PRIMARY KEY,
    recipe_id   INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    category_id INTEGER DEFAULT 1 REFERENCES categories (id) ON DELETE SET DEFAULT,
    UNIQUE (recipe_id, category_id)
);

CREATE TABLE recipe_ingredient
(
    id            INTEGER PRIMARY KEY,
    recipe_id     INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    ingredient_id INTEGER NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, ingredient_id)
);

CREATE TABLE recipe_instruction
(
    id             INTEGER PRIMARY KEY,
    recipe_id      INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    instruction_id INTEGER NOT NULL REFERENCES instructions (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, instruction_id)
);

CREATE TABLE recipe_tool
(
    id        INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    tool_id   INTEGER NOT NULL REFERENCES tools (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, tool_id)
);

CREATE TABLE recipe_keyword
(
    id         INTEGER PRIMARY KEY,
    recipe_id  INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    keyword_id INTEGER NOT NULL REFERENCES keywords (id) ON DELETE CASCADE,
    UNIQUE (recipe_id, keyword_id)
);

CREATE TABLE recipe_time
(
    id        INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    time_id   INTEGER NOT NULL REFERENCES times (id) ON DELETE SET NULL,
    UNIQUE (recipe_id, time_id)
);

CREATE TABLE websites
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    url  TEXT UNIQUE NOT NULL
);

-- +goose StatementBegin
CREATE TRIGGER users_insert_init
    AFTER INSERT
    ON users
    FOR EACH ROW
BEGIN
    INSERT INTO counts (user_id)
    VALUES (NEW.id);

    INSERT INTO user_category (user_id, category_id)
    VALUES (NEW.id, (SELECT id FROM categories WHERE name = 'uncategorized')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'appetizers')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'bread')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'breakfasts')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'condiments')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'dessert')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'lunch')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'main dish')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'salad')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'side dish')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'snacks')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'soups')),
           (NEW.id, (SELECT id FROM categories WHERE name = 'stews'));
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER times_calc_total_time_insert
    AFTER INSERT
    ON times
    FOR EACH ROW
BEGIN
    UPDATE times
    SET total = NEW.prep + NEW.cook
    WHERE id = NEW.id
      AND prep IS NOT NULL
      AND cook IS NOT NULL
      AND total IS NULL;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER times_calc_total_time_update
    AFTER UPDATE
    ON times
    FOR EACH ROW
BEGIN
    UPDATE times
    SET total = NEW.prep + NEW.cook
    WHERE id = NEW.id
      AND prep IS NOT NULL
      AND cook IS NOT NULL
      AND total IS NULL;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER recipes_update_updated_at
    BEFORE UPDATE
    ON recipes
    FOR EACH ROW
BEGIN
    UPDATE recipes
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER user_recipe_insert_inc_count
    AFTER INSERT
    ON user_recipe
    FOR EACH ROW
BEGIN
    UPDATE counts
    SET num_recipes = num_recipes + 1
    WHERE id = NEW.user_id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER user_recipe_delete_dec_count
    AFTER DELETE
    ON user_recipe
    FOR EACH ROW
BEGIN
    UPDATE counts
    SET num_recipes = num_recipes - 1
    WHERE id = OLD.user_id;
END;
-- +goose StatementEnd

-- Inserts
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

INSERT INTO websites (name, url)
VALUES ('101cookbooks.com', 'https://101cookbooks.com'),
       ('claudia.abril.com', 'https://www.claudia.abril.com.br/receitas'),
       ('acouplecooks.com', 'https://www.acouplecooks.com'),
       ('afghankitchenrecipes.com', 'http://www.afghankitchenrecipes.com'),
       ('allrecipes.com', 'https://www.allrecipes.com'),
       ('amazingribs.com', 'https://amazingribs.com'),
       ('ambitiouskitchen.com', 'https://www.ambitiouskitchen.com'),
       ('archanaskitchen.com', 'https://www.archanaskitchen.com'),
       ('atelierdeschefs.fr', 'https://www.atelierdeschefs.fr'),
       ('averiecooks.com', 'https://www.averiecooks.com'),
       ('bakingmischief.com', 'https://bakingmischief.com'),
       ('baking-sense.com', 'https://www.baking-sense.com'),
       ('bbc.co.uk', 'https://www.bbc.co.uk/food/recipes'),
       ('bbcgoodfood.com', 'https://www.bbcgoodfood.com/recipes'),
       ('bettycrocker.com', 'https://www.bettycrocker.com/recipes'),
       ('bigoven.com', 'https://www.bigoven.com'),
       ('bonappetit.com', 'https://www.bonappetit.com'),
       ('bowlofdelicious.com', 'https://www.bowlofdelicious.com'),
       ('budgetbytes.com', 'https://www.budgetbytes.com'),
       ('castironketo.net', 'https://www.castironketo.net'),
       ('cdkitchen.com', 'https://www.cdkitchen.com'),
       ('chefkoch.de', 'https://www.chefkoch.de'),
       ('comidinhasdochef.com', 'https://comidinhasdochef.com'),
       ('cookeatshare.com', 'https://cookeatshare.com'),
       ('cookieandkate.com', 'https://cookieandkate.com'),
       ('cookinglight.com', 'https://www.cookinglight.com'),
       ('cookstr.com', 'https://www.cookstr.com'),
       ('copykat.com', 'https://copykat.com'),
       ('countryliving.com', 'https://www.countryliving.com'),
       ('cuisineaz.com', 'https://www.cuisineaz.com'),
       ('cybercook.com.br', 'https://cybercook.com.br'),
       ('delish.com', 'https://www.delish.com'),
       ('ditchthecarbs.com', 'https://www.ditchthecarbs.com'),
       ('domesticate-me.com', 'https://domesticate-me.com'),
       ('downshiftology.com', 'https://downshiftology.com'),
       ('dr.dk', 'https://www.dr.dk/mad/opskrift'),
       ('eatingbirdfood.com', 'https://www.eatingbirdfood.com'),
       ('eatingwell.com', 'https://www.eatingwell.com'),
       ('eatsmarter.com', 'https://eatsmarter.com'),
       ('eatwhattonight.com', 'https://eatwhattonight.com'),
       ('epicurious.com', 'https://www.epicurious.com'),
       ('expressen.se', 'https://www.expressen.se/alltommat/recept'),
       ('fifteenspatulas.com', 'https://www.fifteenspatulas.com'),
       ('finedininglovers.com', 'https://www.finedininglovers.com'),
       ('fitmencook.com', 'https://fitmencook.com'),
       ('food.com', 'https://www.food.com'),
       ('food52.com', 'https://food52.com/recipes'),
       ('foodandwine.com', 'https://www.foodandwine.com'),
       ('foodrepublic.com', 'https://www.foodrepublic.com'),
       ('forksoverknives.com', 'https://www.forksoverknives.com'),
       ('franzoesischkochen.de', 'https://www.franzoesischkochen.de'),
       ('giallozafferano.com', 'https://www.giallozafferano.com'),
       ('gimmesomeoven.com', 'https://www.gimmesomeoven.com'),
       ('globo.com', 'https://receitas.globo.com'),
       ('gonnawantseconds.com', 'https://www.gonnawantseconds.com'),
       ('greatbritishchefs.com', 'https://www.greatbritishchefs.com'),
       ('halfbakedharvest.com', 'https://www.halfbakedharvest.com'),
       ('hassanchef.com', 'https://www.hassanchef.com'),
       ('headbangerskitchen.com', 'https://headbangerskitchen.com'),
       ('hellofresh.com', 'https://www.hellofresh.com/recipes'),
       ('homechef.com', 'https://www.homechef.com/our-menu'),
       ('hostthetoast.com', 'https://hostthetoast.com'),
       ('indianhealthyrecipes.com', 'https://www.indianhealthyrecipes.com'),
       ('innit.com', 'https://www.innit.com/meal'),
       ('inspiralized.com', 'https://inspiralized.com'),
       ('jamieoliver.html.com', 'https://www.jamieoliver.html.com'),
       ('jimcooksfoodgood.com', 'https://jimcooksfoodgood.com'),
       ('joyfoodsunshine.com', 'https://joyfoodsunshine.com'),
       ('justataste.com', 'https://www.justataste.com'),
       ('justonecookbook.com', 'https://www.justonecookbook.com'),
       ('kennymcgovern.com', 'https://kennymcgovern.com'),
       ('kingarthurbaking.com', 'https://www.kingarthurbaking.com'),
       ('kochbar.de', 'https://www.kochbar.de/rezept'),
       ('koket.se', 'https://www.koket.se'),
       ('kuchnia-domowa.pl', 'https://www.kuchnia-domowa.pl'),
       ('kwestiasmaku.com', 'https://www.kwestiasmaku.com'),
       ('lecremedelacrumb.com', 'https://www.lecremedelacrumb.com'),
       ('lekkerensimpel.com', 'https://www.lekkerensimpel.com'),
       ('littlespicejar.com', 'https://littlespicejar.com'),
       ('livelytable.com', 'https://livelytable.com'),
       ('lovingitvegan.com', 'https://lovingitvegan.com'),
       ('madensverden.dk', 'https://madensverden.dk'),
       ('marthastewart.com', 'https://www.marthastewart.com'),
       ('matprat.no', 'https://www.matprat.no'),
       ('melskitchencafe.com', 'https://www.melskitchencafe.com'),
       ('mindmegette.hu', 'https://www.mindmegette.hu'),
       ('minimalistbaker.com', 'https://minimalistbaker.com'),
       ('misya.info', 'https://www.misya.info'),
       ('momswithcrockpots.com', 'https://momswithcrockpots.com'),
       ('monsieur-cuisine.com', 'https://www.monsieur-cuisine.com'),
       ('motherthyme.com', 'https://www.motherthyme.com'),
       ('mybakingaddiction.com', 'https://www.mybakingaddiction.com'),
       ('mykitchen101.com', 'https://mykitchen101.com'),
       ('mykitchen101en.com', 'https://mykitchen101en.com'),
       ('myplate.gov', 'https://www.myplate.gov/recipes'),
       ('myrecipes.com', 'https://www.myrecipes.com'),
       ('nourishedbynutrition.com', 'https://nourishedbynutrition.com'),
       ('nutritionbynathalie.com', 'https://www.nutritionbynathalie.com'),
       ('nytimes.com', 'https://cooking.nytimes.com/recipes'),
       ('ohsheglows.com', 'https://ohsheglows.com'),
       ('onceuponachef.com', 'https://www.onceuponachef.com'),
       ('paleorunningmomma.com', 'https://www.paleorunningmomma.com'),
       ('panelinha.com.br', 'https://www.panelinha.com.br'),
       ('paninihappy.com', 'https://paninihappy.com'),
       ('practicalselfreliance.com', 'https://practicalselfreliance.com'),
       ('primaledgehealth.com', 'https://www.primaledgehealth.com'),
       ('przepisy.pl', 'https://www.przepisy.pl'),
       ('purelypope.com', 'https://purelypope.com'),
       ('purplecarrot.com', 'https://www.purplecarrot.com'),
       ('rachlmansfield.com', 'https://rachlmansfield.com'),
       ('rainbowplantlife.com', 'https://rainbowplantlife.com'),
       ('realsimple.com', 'https://www.realsimple.com'),
       ('recipetineats.com', 'https://www.recipetineats.com'),
       ('redhousespice.com', 'https://redhousespice.com'),
       ('reishunger.de', 'https://www.reishunger.de'),
       ('rezeptwelt.de', 'https://www.rezeptwelt.de'),
       ('sallysbakingaddiction.com', 'https://sallysbakingaddiction.com'),
       ('saveur.com', 'https://www.saveur.com'),
       ('seriouseats.com', 'https://www.seriouseats.com'),
       ('simplyquinoa.com', 'https://www.simplyquinoa.com'),
       ('simplyrecipes.com', 'https://www.simplyrecipes.com'),
       ('simplywhisked.com', 'https://www.simplywhisked.com'),
       ('skinnytaste.com', 'https://www.skinnytaste.com'),
       ('southernliving.com', 'https://www.southernliving.com'),
       ('spendwithpennies.com', 'https://www.spendwithpennies.com'),
       ('steamykitchen.com', 'https://steamykitchen.com'),
       ('streetkitchen.co', 'https://streetkitchen.co'),
       ('sunbasket.com', 'https://sunbasket.com'),
       ('sweetcsdesigns.com', 'https://sweetcsdesigns.com'),
       ('sweetpeasandsaffron.com', 'https://sweetpeasandsaffron.com'),
       ('tasteofhome.com', 'https://www.tasteofhome.com'),
       ('tastesoflizzyt.com', 'https://www.tastesoflizzyt.com'),
       ('tasty.co', 'https://tasty.co'),
       ('tastykitchen.com', 'https://tastykitchen.com'),
       ('tesco.com', 'https://realfood.tesco.com'),
       ('theclevercarrot.com', 'https://www.theclevercarrot.com'),
       ('thehappyfoodie.co.uk', 'https://thehappyfoodie.co.uk'),
       ('thekitchenmagpie.com', 'https://www.thekitchenmagpie.com'),
       ('thenutritiouskitchen.co', 'http://thenutritiouskitchen.co'),
       ('thepioneerwoman.com', 'https://www.thepioneerwoman.com'),
       ('thespruceeats.com', 'https://www.thespruceeats.com'),
       ('thevintagemixer.com', 'https://www.thevintagemixer.com'),
       ('thewoksoflife.com', 'https://thewoksoflife.com'),
       ('timesofindia.com', 'https://recipes.timesofindia.com/recipes'),
       ('tine.no', 'https://www.tine.no/oppskrifter'),
       ('twopeasandtheirpod.com', 'https://www.twopeasandtheirpod.com'),
       ('valdemarsro.dk', 'https://www.valdemarsro.dk'),
       ('vanillaandbean.com', 'https://vanillaandbean.com'),
       ('vegolosi.it', 'https://www.vegolosi.it'),
       ('vegrecipesofindia.com', 'https://www.vegrecipesofindia.com'),
       ('watchwhatueat.com', 'https://www.watchwhatueat.com'),
       ('whatsgabycooking.com', 'https://whatsgabycooking.com'),
       ('wikibooks.org', 'https://en.wikibooks.org'),
       ('m.wikibooks.org', 'https://en.m.wikibooks.org'),
       ('woop.co.nz', 'https://woop.co.nz'),
       ('ye-mek.net', 'https://ye-mek.net'),
       ('zenbelly.com', 'https://www.zenbelly.com');

-- +goose Down
DROP TABLE IF EXISTS user_recipe;
DROP TABLE IF EXISTS user_category;
DROP TABLE IF EXISTS recipe_category;
DROP TABLE IF EXISTS recipe_ingredient;
DROP TABLE IF EXISTS recipe_instruction;
DROP TABLE IF EXISTS recipe_tool;
DROP TABLE IF EXISTS recipe_keyword;
DROP TABLE IF EXISTS recipe_time;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS instructions;
DROP TABLE IF EXISTS tools;
DROP TABLE IF EXISTS keywords;
DROP TABLE IF EXISTS nutrition;
DROP TABLE IF EXISTS times;
DROP TABLE IF EXISTS counts;
DROP TABLE IF EXISTS fruits_veggies;
DROP TABLE IF EXISTS blacklist_units;
DROP TABLE IF EXISTS websites;

DROP TRIGGER IF EXISTS users_insert_init;
DROP TRIGGER IF EXISTS times_calc_total_time_insert;
DROP TRIGGER IF EXISTS times_calc_total_time_update;
DROP TRIGGER IF EXISTS recipes_update_updated_at;
DROP TRIGGER IF EXISTS user_recipe_insert_inc_count;
DROP TRIGGER IF EXISTS user_recipe_delete_dec_count;