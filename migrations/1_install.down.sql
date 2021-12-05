DROP EXTENSION "uuid-ossp";
DROP FUNCTION times_calc_total_time() CASCADE;
DROP FUNCTION recipes_update_updated_at() CASCADE;

TRUNCATE TABLE recipes CASCADE;

DROP TABLE recipes CASCADE;
DROP TABLE categories CASCADE;
DROP TABLE ingredients CASCADE;
DROP TABLE instructions CASCADE;
DROP TABLE nutrition;
DROP TABLE tools CASCADE;
DROP TABLE keywords CASCADE;
DROP TABLE times CASCADE;
DROP TABLE category_recipe;
DROP TABLE ingredient_recipe;
DROP TABLE instruction_recipe;
DROP TABLE time_recipe;
DROP TABLE tool_recipe;
DROP TABLE keyword_recipe;
