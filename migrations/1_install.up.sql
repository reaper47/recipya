--
-- CREATES
--
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE times (
  id SERIAL PRIMARY KEY,
  prep INTERVAL,
  cook INTERVAL,
  total INTERVAL,
  UNIQUE (prep, cook)
);

CREATE TABLE nutrition (
  id SERIAL PRIMARY KEY,
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

CREATE TABLE recipes (
  id SERIAL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  description TEXT,
  url VARCHAR(80),
  image UUID,
  yield SMALLINT,
  category_id INTEGER DEFAULT 0 REFERENCES categories(id) ON DELETE SET DEFAULT,
  times_id INTEGER REFERENCES times(id) ON DELETE SET NULL,
  nutrition_id INTEGER REFERENCES nutrition(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ingredients (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE instructions (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE tools (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE keywords (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
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

--
-- TRIGGERS
--
CREATE FUNCTION times_calc_total_time() RETURNS trigger AS $func$
  BEGIN
    IF NEW.prep IS NOT NULL AND NEW.cook IS NOT NULL AND NEW.total IS NULL THEN
      NEW.total := NEW.prep + NEW.cook;
      RETURN NEW;
    END IF;
  END;
$func$ LANGUAGE plpgsql;

CREATE TRIGGER times_calc_total_time 
  BEFORE INSERT OR UPDATE ON times 
  FOR EACH ROW 
  EXECUTE FUNCTION times_calc_total_time();

--
-- INSERTS
--
INSERT INTO categories (name) 
VALUES ('uncategorized');
