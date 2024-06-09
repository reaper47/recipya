-- +goose Up
INSERT INTO websites (host, url)
VALUES
    ('chatelaine.com', 'https://chatelaine.com/recipe/salads/mozzarella-clementine-panzanella/?utm_source=whisk&utm_medium=webapp&utm_campaign=fresh_mozzarella_and_clementine_panzanella'),
    ('downshiftology.com', 'https://downshiftology.com/recipes/shakshuka/?utm_source=whisk&utm_medium=webapp&utm_campaign=shakshuka_recipe_(easy_%26_traditional)'),
    ('natashaskitchen.com', 'https://natashaskitchen.com/taco-seasoning-recipe/'),
    ('wellplated.com', 'https://www.wellplated.com/kale-pineapple-smoothie/')
;

-- +goose Down
DELETE FROM websites
WHERE (host = 'chatelaine.com/' AND url = 'https://chatelaine.com/recipe/salads/mozzarella-clementine-panzanella/?utm_source=whisk&utm_medium=webapp&utm_campaign=fresh_mozzarella_and_clementine_panzanella') OR
      (host = 'downshiftology.com' AND url = 'https://downshiftology.com/recipes/shakshuka/?utm_source=whisk&utm_medium=webapp&utm_campaign=shakshuka_recipe_(easy_%26_traditional)') OR
      (host = 'natashaskitchen.com/' AND url = 'https://natashaskitchen.com/taco-seasoning-recipe/') OR
      (host = 'wellplated.com' AND url = 'https://www.wellplated.com/kale-pineapple-smoothie/')
