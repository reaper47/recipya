-- +goose Up
--
-- Create normal tables
--
CREATE TABLE recipes
(
    id          INTEGER PRIMARY KEY,
    name        TEXT      NOT NULL,
    description TEXT,
    image       TEXT               DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' ||
                                            substr(lower(hex(randomblob(2))), 2) || '-a' ||
                                            substr(lower(hex(randomblob(2))), 2) || '-%' ||
                                            substr(lower(hex(randomblob(6))), 2)),
    yield       INTEGER,
    url         TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
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
    id               INTEGER PRIMARY KEY,
    ingredient_id    INTEGER NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    recipe_id        INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    ingredient_order INTEGER NOT NULL
);

CREATE TABLE instruction_recipe
(
    id                INTEGER PRIMARY KEY,
    instruction_id    INTEGER NOT NULL REFERENCES instructions (id) ON DELETE CASCADE,
    recipe_id         INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    instruction_order INTEGER NOT NULL
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
