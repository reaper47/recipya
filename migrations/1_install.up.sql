--
-- CREATES
--
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

CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE ingredients (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE instructions (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE tools (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE keywords (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
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
  recipes INTEGER NOT NULL
);


--
-- Association Tables
--
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

CREATE FUNCTION recipes_insert_inc_count() RETURNS TRIGGER AS
$BODY$
BEGIN 
  UPDATE counts SET 
    recipes = recipes + 1
  WHERE id = 1;
  RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

CREATE FUNCTION recipes_delete_dec_count() RETURNS TRIGGER AS
$BODY$
BEGIN 
  UPDATE counts SET 
    recipes = recipes - 1
  WHERE id = 1;
  RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;


--
-- TRIGGERS
--
CREATE TRIGGER times_calc_total_time 
  BEFORE INSERT OR UPDATE ON times 
  FOR EACH ROW 
  EXECUTE FUNCTION times_calc_total_time();

CREATE TRIGGER recipes_update_updated_at
	BEFORE UPDATE ON recipes 
	FOR EACH ROW 
	EXECUTE FUNCTION recipes_update_updated_at();

CREATE TRIGGER recipes_insert_inc_count
  AFTER INSERT ON recipes
  FOR EACH ROW
  EXECUTE PROCEDURE recipes_insert_inc_count();

CREATE TRIGGER recipes_delete_dec_count
  AFTER DELETE ON recipes
  FOR EACH ROW
  EXECUTE PROCEDURE recipes_delete_dec_count();

--
-- INSERTS
--
INSERT INTO categories (name) 
VALUES ('uncategorized');

INSERT INTO counts (recipes)
VALUES (0);

--
-- EXTENSIONS
--
CREATE EXTENSION "uuid-ossp";
