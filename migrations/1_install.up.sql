--
-- EXTENSIONS
--
CREATE EXTENSION pgcrypto;
CREATE EXTENSION "uuid-ossp";

--
-- CREATES
--
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(36) NOT NULL UNIQUE CHECK (LOWER(username) = username),
	email VARCHAR(254) NOT NULL UNIQUE CHECK (LOWER(email) = email),
	hashed_password VARCHAR(254) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipes (
	id SERIAL PRIMARY KEY,
	name VARCHAR(80) NOT NULL,
	description TEXT,
	url VARCHAR(256),
	image UUID DEFAULT gen_random_uuid(),
	yield SMALLINT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- https://recipeland.com/recipes/categories/browse
CREATE TABLE categories (
	id SERIAL PRIMARY KEY,
	name VARCHAR(24) UNIQUE NOT NULL CHECK (LOWER(name) = name)
);

CREATE TABLE ingredients (
	id SERIAL PRIMARY KEY,
	name VARCHAR(24) UNIQUE NOT NULL CHECK (LOWER(name) = name)
);

CREATE TABLE instructions (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL CHECK (LOWER(name) = name)
);

CREATE TABLE tools (
	id SERIAL PRIMARY KEY,
	name VARCHAR(24) UNIQUE NOT NULL CHECK (LOWER(name) = name)
);

CREATE TABLE keywords (
	id SERIAL PRIMARY KEY,
	name VARCHAR(24) UNIQUE NOT NULL CHECK (LOWER(name) = name)
);

CREATE TABLE nutrition (
	id SERIAL PRIMARY KEY,
	recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
	calories VARCHAR(10),
	total_carbohydrates VARCHAR(5),
	sugars VARCHAR(7),
	protein VARCHAR(5),
	total_fat VARCHAR(5),
	saturated_fat VARCHAR(5),
	cholesterol VARCHAR(7),
	sodium VARCHAR(7),
	fiber VARCHAR(5)
);

CREATE TABLE times (
	id SERIAL PRIMARY KEY,
	prep INTERVAL,
	cook INTERVAL,
	total INTERVAL,
	UNIQUE (prep, cook)
);

CREATE TABLE counts (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	recipes INTEGER DEFAULT 0
);

--
-- Association Tables
--
CREATE TABLE user_recipe (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
	UNIQUE (user_id, recipe_id)
);

CREATE TABLE user_category (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
	UNIQUE (user_id, category_id)
);

CREATE TABLE category_recipe (
	id SERIAL PRIMARY KEY,
	recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
	category_id INTEGER DEFAULT 1 REFERENCES categories(id) ON DELETE SET DEFAULT,
	UNIQUE (recipe_id, category_id)
);

CREATE TABLE ingredient_recipe (
	id SERIAL PRIMARY KEY,
	recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
	ingredient_id INTEGER NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
	UNIQUE (recipe_id, ingredient_id)
);

CREATE TABLE instruction_recipe (
	id SERIAL PRIMARY KEY,
	recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
	instruction_id INTEGER NOT NULL REFERENCES instructions(id) ON DELETE CASCADE,
	UNIQUE (recipe_id, instruction_id)
);

CREATE TABLE tool_recipe (
	id SERIAL PRIMARY KEY,
	recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
	tool_id INTEGER NOT NULL REFERENCES tools(id) ON DELETE CASCADE,
	UNIQUE (recipe_id, tool_id)
);

CREATE TABLE keyword_recipe (
	id SERIAL PRIMARY KEY,
	recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
	keyword_id INTEGER NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
	UNIQUE (recipe_id, keyword_id)
);

CREATE TABLE time_recipe (
	id SERIAL PRIMARY KEY,
	recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
	time_id INTEGER NOT NULL REFERENCES times(id) ON DELETE SET NULL,
	UNIQUE (recipe_id, time_id)
);

--
-- FUNCTIONS
--
CREATE FUNCTION users_insert_init() RETURNS TRIGGER AS
$BODY$
BEGIN
	INSERT INTO counts (user_id)
	VALUES (NEW.id);

	INSERT INTO user_category (user_id, category_id)
	VALUES
		(NEW.id, (SELECT id FROM categories WHERE name = 'uncategorized')),
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

	RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

CREATE FUNCTION times_calc_total_time() RETURNS TRIGGER AS
$BODY$
BEGIN
	IF NEW.prep IS NOT NULL AND NEW.cook IS NOT NULL AND NEW.total IS NULL THEN
		NEW.total := NEW.prep + NEW.cook;
		RETURN NEW;
	END IF;
END;
$BODY$
LANGUAGE plpgsql;

CREATE FUNCTION recipes_update_updated_at() RETURNS TRIGGER AS
$BODY$
BEGIN
	NEW.updated_at := CURRENT_TIMESTAMP;
	RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

CREATE FUNCTION user_recipe_insert_inc_count() RETURNS TRIGGER AS
$BODY$
BEGIN
	UPDATE counts
	SET recipes = recipes + 1
	WHERE id = NEW.user_id;

	RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

CREATE FUNCTION user_recipe_delete_dec_count() RETURNS TRIGGER AS
$BODY$
BEGIN
	UPDATE counts
	SET recipes = recipes - 1
	WHERE id = OLD.user_id;
	RETURN OLD;
END;
$BODY$
LANGUAGE plpgsql;

--
-- TRIGGERS
--
CREATE TRIGGER users_insert_init
  AFTER INSERT ON users
  FOR EACH ROW
  EXECUTE FUNCTION users_insert_init();

CREATE TRIGGER times_calc_total_time
  BEFORE INSERT OR UPDATE ON times
  FOR EACH ROW
  EXECUTE FUNCTION times_calc_total_time();

CREATE TRIGGER recipes_update_updated_at
	BEFORE UPDATE ON recipes
	FOR EACH ROW
	EXECUTE FUNCTION recipes_update_updated_at();

CREATE TRIGGER user_recipe_insert_inc_count
  AFTER INSERT ON user_recipe
  FOR EACH ROW
  EXECUTE PROCEDURE user_recipe_insert_inc_count();

CREATE TRIGGER user_recipe_delete_dec_count
  AFTER DELETE ON user_recipe
  FOR EACH ROW
  EXECUTE PROCEDURE user_recipe_delete_dec_count();

--
-- INSERTS
--
INSERT INTO categories (name)
VALUES
	('uncategorized'),
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
