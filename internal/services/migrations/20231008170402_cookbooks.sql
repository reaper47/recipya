-- +goose Up
CREATE TABLE cookbooks
(
    id      INTEGER PRIMARY KEY,
    title   TEXT NOT NULL,
    image   TEXT,
    count   INTEGER DEFAULT 0,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (title, user_id)
);

CREATE TABLE cookbook_recipes
(
    id          INTEGER PRIMARY KEY,
    cookbook_id INTEGER REFERENCES cookbooks (id) ON DELETE CASCADE,
    recipe_id   INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL,
    UNIQUE (cookbook_id, recipe_id)
);

CREATE VIRTUAL TABLE cookbooks_fts USING fts5
(
    id,
    user_id,
    title
);

CREATE VIRTUAL TABLE recipes_fts USING fts5
(
    id,
    user_id,
    name,
    description,
    category,
    cuisine,
    ingredients,
    instructions,
    keywords,
    source
);

ALTER TABLE user_settings
    ADD COLUMN cookbooks_view INTEGER DEFAULT 0;

ALTER TABLE counts
    ADD COLUMN cookbooks INTEGER DEFAULT 0;

CREATE TABLE shadow_last_inserted_recipe
(
    row         INTEGER PRIMARY KEY,
    id          INTEGER NOT NULL,
    name        TEXT    NOT NULL,
    description TEXT,
    source      TEXT
);

-- +goose StatementBegin
CREATE TRIGGER cookbook_recipes_insert
    AFTER INSERT
    ON cookbook_recipes
    FOR EACH ROW
BEGIN
    UPDATE cookbooks
    SET count = count + 1
    WHERE NEW.cookbook_id = id;
END;

CREATE TRIGGER cookbook_recipes_delete
    AFTER DELETE
    ON cookbook_recipes
    FOR EACH ROW
BEGIN
    UPDATE cookbooks
    SET count = count - 1
    WHERE OLD.cookbook_id = id;
END;

CREATE TRIGGER cookbooks_insert
    AFTER INSERT
    ON cookbooks
    FOR EACH ROW
BEGIN
    UPDATE counts
    SET cookbooks = cookbooks + 1
    WHERE user_id = NEW.user_id;

    INSERT INTO cookbooks_fts (id, user_id, title)
    VALUES (NEW.id, NEW.user_id, NEW.title);
END;

CREATE TRIGGER cookbooks_delete
    AFTER DELETE
    ON cookbooks
    FOR EACH ROW
BEGIN
    UPDATE counts
    SET cookbooks = cookbooks - 1
    WHERE user_id = OLD.user_id;

    DELETE
    FROM cookbooks_fts
    WHERE id = OLD.id
      AND user_id = OLD.user_id;
END;

CREATE TRIGGER shadow_last_inserted_recipe_insert
    AFTER INSERT
    ON shadow_last_inserted_recipe
    FOR EACH ROW
BEGIN
    INSERT INTO recipes_fts (id,
                             user_id,
                             name,
                             description,
                             category,
                             cuisine,
                             ingredients,
                             instructions,
                             keywords,
                             source)
    VALUES (NEW.id,
            (SELECT user_id FROM user_recipe AS ur WHERE ur.recipe_id = NEW.id),
            NEW.name,
            NEW.description,
            (SELECT c.name
             FROM category_recipe AS cr
                      JOIN categories AS c ON cr.category_id = c.id
             WHERE cr.recipe_id = NEW.id),
            (SELECT c.name
             FROM cuisine_recipe AS cr
                      JOIN categories AS c ON cr.cuisine_id = c.id
             WHERE cr.recipe_id = NEW.id),
            (SELECT COALESCE((SELECT GROUP_CONCAT(ingredient_name, '<!---->')
                              FROM (SELECT DISTINCT ingredients.name AS ingredient_name
                                    FROM ingredient_recipe
                                             JOIN ingredients ON ingredients.id = ingredient_recipe.ingredient_id
                                    WHERE ingredient_recipe.recipe_id = NEW.id
                                    ORDER BY ingredient_order)), '')),
            (SELECT COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
                              FROM (SELECT DISTINCT instructions.name AS instruction_name
                                    FROM instruction_recipe
                                             JOIN instructions ON instructions.id = instruction_recipe.instruction_id
                                    WHERE instruction_recipe.recipe_id = NEW.id
                                    ORDER BY instruction_order)), '')),
            (SELECT COALESCE((SELECT GROUP_CONCAT(keyword_name, ',')
                              FROM (SELECT DISTINCT keywords.name AS keyword_name
                                    FROM keyword_recipe
                                             JOIN keywords ON keywords.id = keyword_recipe.keyword_id
                                    WHERE keyword_recipe.recipe_id = NEW.id)), '')),
            NEW.source);
END;

CREATE TRIGGER shadow_last_inserted_recipe_delete
    AFTER DELETE
    ON shadow_last_inserted_recipe
    FOR EACH ROW
BEGIN
    DELETE FROM recipes_fts WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER cookbook_recipes_insert;
DROP TRIGGER cookbook_recipes_delete;
DROP TRIGGER cookbooks_insert;
DROP TRIGGER cookbooks_delete;
DROP TRIGGER shadow_last_inserted_recipe_insert;
DROP TRIGGER shadow_last_inserted_recipe_delete;
ALTER TABLE user_settings
    DROP COLUMN cookbooks_view;
ALTER TABLE counts
    DROP COLUMN cookbooks;
DROP TABLE cookbook_recipes;
DROP TABLE cookbooks;
DROP TABLE cookbooks_fts;
