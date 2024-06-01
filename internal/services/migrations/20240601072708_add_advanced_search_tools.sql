-- +goose Up
-- +goose StatementBegin
DROP TABLE recipes_fts;
DROP TRIGGER shadow_last_inserted_recipe_insert;
DROP TRIGGER trig_update_recipe_buo;

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
    tools,
    source
);

INSERT INTO recipes_fts (id,
                         user_id,
                         name,
                         description,
                         category,
                         cuisine,
                         ingredients,
                         instructions,
                         keywords,
                         tools,
                         source)
SELECT recipes.id,
       ur.user_id,
       recipes.name,
       recipes.description,
       categories.name,
       cuisines.name,
       COALESCE((SELECT GROUP_CONCAT(ingredient_name, '<!---->')
                 FROM (SELECT DISTINCT ingredients.name AS ingredient_name
                       FROM ingredient_recipe
                                JOIN ingredients ON ingredients.id = ingredient_recipe.ingredient_id
                       WHERE ingredient_recipe.recipe_id = recipes.id
                       ORDER BY ingredient_order)), ''),
       COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
                 FROM (SELECT DISTINCT instructions.name AS instruction_name
                       FROM instruction_recipe
                                JOIN instructions ON instructions.id = instruction_recipe.instruction_id
                       WHERE instruction_recipe.recipe_id = recipes.id
                       ORDER BY instruction_order)), ''),
       GROUP_CONCAT(DISTINCT keywords.name) AS keywords,
       (SELECT GROUP_CONCAT(name)
        FROM (SELECT tool_recipe.quantity || ' ' || tools.name AS name
              FROM tool_recipe
                       JOIN tools ON tool_recipe.tool_id = tools.id
              WHERE tool_recipe.recipe_id = recipes.id
              ORDER BY tool_recipe.tool_order)),
       recipes.url                          AS source
FROM recipes
         LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id
         LEFT JOIN categories ON category_recipe.category_id = categories.id
         LEFT JOIN cuisine_recipe ON recipes.id = cuisine_recipe.recipe_id
         LEFT JOIN cuisines ON cuisine_recipe.cuisine_id = cuisines.id
         LEFT JOIN ingredient_recipe ON recipes.id = ingredient_recipe.recipe_id
         LEFT JOIN ingredients ON ingredient_recipe.ingredient_id = ingredients.id
         LEFT JOIN instruction_recipe ON recipes.id = instruction_recipe.recipe_id
         LEFT JOIN instructions ON instruction_recipe.instruction_id = instructions.id
         LEFT JOIN keyword_recipe ON recipes.id = keyword_recipe.recipe_id
         LEFT JOIN keywords ON keyword_recipe.keyword_id = keywords.id
         LEFT JOIN tool_recipe ON recipes.id = tool_recipe.recipe_id
         LEFT JOIN tools ON tool_recipe.tool_id = tools.id
         LEFT JOIN nutrition ON recipes.id = nutrition.recipe_id
         LEFT JOIN time_recipe ON recipes.id = time_recipe.recipe_id
         LEFT JOIN times ON time_recipe.time_id = times.id
         INNER JOIN user_recipe AS ur ON ur.recipe_id = recipes.id
GROUP BY recipes.id;

CREATE TRIGGER trig_shadow_last_inserted_recipe_ai
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
                             tools,
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
            (SELECT GROUP_CONCAT(name)
             FROM (SELECT tool_recipe.quantity || ' ' || tools.name AS name
                   FROM tool_recipe
                            JOIN tools ON tool_recipe.tool_id = tools.id
                   WHERE tool_recipe.recipe_id = NEW.id
                   ORDER BY tool_recipe.tool_order)),
            NEW.source);
END;

CREATE TRIGGER trig_update_recipe_buo
    BEFORE UPDATE OF id
    ON recipes
    FOR EACH ROW
BEGIN
    UPDATE recipes_fts
    SET name         = NEW.name,
        description  = NEW.description,
        category     = (SELECT c.name
                        FROM category_recipe AS cr
                                 JOIN categories AS c ON cr.category_id = c.id
                        WHERE cr.recipe_id = NEW.id),
        cuisine      = (SELECT c.name
                        FROM cuisine_recipe AS cr
                                 JOIN cuisines c ON cr.cuisine_id = c.id
                        WHERE cr.recipe_id = NEW.id),
        ingredients  = (SELECT COALESCE((SELECT GROUP_CONCAT(ingredient_name, '<!---->')
                                         FROM (SELECT DISTINCT ingredients.name AS ingredient_name
                                               FROM ingredient_recipe
                                                        JOIN ingredients ON ingredients.id = ingredient_recipe.ingredient_id
                                               WHERE ingredient_recipe.recipe_id = NEW.id
                                               ORDER BY ingredient_order)), '')),
        instructions = (SELECT COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
                                         FROM (SELECT DISTINCT instructions.name AS instruction_name
                                               FROM instruction_recipe
                                                        JOIN instructions ON instructions.id = instruction_recipe.instruction_id
                                               WHERE instruction_recipe.recipe_id = NEW.id
                                               ORDER BY instruction_order)), '')),
        keywords     = (SELECT COALESCE((SELECT GROUP_CONCAT(keyword_name, ',')
                                         FROM (SELECT DISTINCT keywords.name AS keyword_name
                                               FROM keyword_recipe
                                                        JOIN keywords ON keywords.id = keyword_recipe.keyword_id
                                               WHERE keyword_recipe.recipe_id = NEW.id)), '')),
        tools        = (SELECT GROUP_CONCAT(name)
                        FROM (SELECT tool_recipe.quantity || ' ' || tools.name AS name
                              FROM tool_recipe
                                       JOIN tools ON tool_recipe.tool_id = tools.id
                              WHERE tool_recipe.recipe_id = NEW.id
                              ORDER BY tool_recipe.tool_order)),
        source       = NEW.url
    WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE recipes_fts;

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
SELECT recipes.id,
       ur.user_id,
       recipes.name,
       recipes.description,
       categories.name,
       cuisines.name,
       COALESCE((SELECT GROUP_CONCAT(ingredient_name, '<!---->')
                 FROM (SELECT DISTINCT ingredients.name AS ingredient_name
                       FROM ingredient_recipe
                                JOIN ingredients ON ingredients.id = ingredient_recipe.ingredient_id
                       WHERE ingredient_recipe.recipe_id = recipes.id
                       ORDER BY ingredient_order)), ''),
       COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
                 FROM (SELECT DISTINCT instructions.name AS instruction_name
                       FROM instruction_recipe
                                JOIN instructions ON instructions.id = instruction_recipe.instruction_id
                       WHERE instruction_recipe.recipe_id = recipes.id
                       ORDER BY instruction_order)), ''),
       GROUP_CONCAT(DISTINCT keywords.name) AS keywords,
       recipes.url                          AS source
FROM recipes
         LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id
         LEFT JOIN categories ON category_recipe.category_id = categories.id
         LEFT JOIN cuisine_recipe ON recipes.id = cuisine_recipe.recipe_id
         LEFT JOIN cuisines ON cuisine_recipe.cuisine_id = cuisines.id
         LEFT JOIN ingredient_recipe ON recipes.id = ingredient_recipe.recipe_id
         LEFT JOIN ingredients ON ingredient_recipe.ingredient_id = ingredients.id
         LEFT JOIN instruction_recipe ON recipes.id = instruction_recipe.recipe_id
         LEFT JOIN instructions ON instruction_recipe.instruction_id = instructions.id
         LEFT JOIN keyword_recipe ON recipes.id = keyword_recipe.recipe_id
         LEFT JOIN keywords ON keyword_recipe.keyword_id = keywords.id
         LEFT JOIN tool_recipe ON recipes.id = tool_recipe.recipe_id
         LEFT JOIN tools ON tool_recipe.tool_id = tools.id
         LEFT JOIN nutrition ON recipes.id = nutrition.recipe_id
         LEFT JOIN time_recipe ON recipes.id = time_recipe.recipe_id
         LEFT JOIN times ON time_recipe.time_id = times.id
         INNER JOIN user_recipe AS ur ON ur.recipe_id = recipes.id
GROUP BY recipes.id;

DROP TRIGGER trig_shadow_last_inserted_recipe_ai;

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

DROP TRIGGER trig_update_recipe_buo;

CREATE TRIGGER trig_update_recipe_buo
    BEFORE UPDATE OF id
    ON recipes
    FOR EACH ROW
BEGIN
    UPDATE recipes_fts
    SET name         = NEW.name,
        description  = NEW.description,
        category     = (SELECT c.name
                        FROM category_recipe AS cr
                                 JOIN categories AS c ON cr.category_id = c.id
                        WHERE cr.recipe_id = NEW.id),
        cuisine      = (SELECT c.name
                        FROM cuisine_recipe AS cr
                                 JOIN cuisines c ON cr.cuisine_id = c.id
                        WHERE cr.recipe_id = NEW.id),
        ingredients  = (SELECT COALESCE((SELECT GROUP_CONCAT(ingredient_name, '<!---->')
                                         FROM (SELECT DISTINCT ingredients.name AS ingredient_name
                                               FROM ingredient_recipe
                                                        JOIN ingredients ON ingredients.id = ingredient_recipe.ingredient_id
                                               WHERE ingredient_recipe.recipe_id = NEW.id
                                               ORDER BY ingredient_order)), '')),
        instructions = (SELECT COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
                                         FROM (SELECT DISTINCT instructions.name AS instruction_name
                                               FROM instruction_recipe
                                                        JOIN instructions ON instructions.id = instruction_recipe.instruction_id
                                               WHERE instruction_recipe.recipe_id = NEW.id
                                               ORDER BY instruction_order)), '')),
        keywords     = (SELECT COALESCE((SELECT GROUP_CONCAT(keyword_name, ',')
                                         FROM (SELECT DISTINCT keywords.name AS keyword_name
                                               FROM keyword_recipe
                                                        JOIN keywords ON keywords.id = keyword_recipe.keyword_id
                                               WHERE keyword_recipe.recipe_id = NEW.id)), '')),
        source       = NEW.url
    WHERE id = NEW.id;
END;
-- +goose StatementEnd
