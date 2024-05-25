-- +goose Up
-- +goose StatementBegin
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
                                 JOIN cuisine AS c ON cr.cuisine_id = c.id
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

-- +goose Down
DROP TRIGGER trig_update_recipe_buo;