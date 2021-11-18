--
-- CREATES
--
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(24) UNIQUE NOT NULL
);

CREATE TABLE times (
  id SERIAL PRIMARY KEY,
  prep_time INTERVAL,
  cook_time INTERVAL,
  total_time INTERVAL
);

CREATE TABLE nutrition (
  id SERIAL PRIMARY KEY,
  calories VARCHAR(10),
  total_carbohydrate VARCHAR(5),
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
  ingredient_id INTEGER NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE
);

CREATE TABLE instruction_recipe (
  id SERIAL PRIMARY KEY,
  recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
  instruction_id INTEGER NOT NULL REFERENCES instructions(id) ON DELETE CASCADE

);

CREATE TABLE tool_recipe (
  id SERIAL PRIMARY KEY,
  recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
  tool_id INTEGER NOT NULL REFERENCES tools(id) ON DELETE CASCADE
);

CREATE TABLE keyword_recipe (
  id SERIAL PRIMARY KEY,
  recipe_id INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
  keyword_id INTEGER NOT NULL REFERENCES keywords(id) ON DELETE CASCADE
);

--
-- TRIGGERS
--
CREATE FUNCTION times_calc_total_time() RETURNS trigger AS $func$
  BEGIN
    IF NEW.prep_time IS NOT NULL AND cook_time IS NOT NULL AND total_time IS NULL THEN
      NEW.total_time := NEW.prep_time + NEW.cook_time;
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
