-- +goose Up
DELETE
FROM websites
WHERE host IN (
               'thenutritiouskitchen.co', -- Removed because the website doesn't work anymore.
               'foodbag.be' -- Removed because the API URL doesn't work anymore. The data comes from somewhere else.
    );


-- +goose Down
INSERT INTO websites (host, url)
VALUES ('thenutritiouskitchen.co', 'https:///thenutritiouskitchen.co'),
       ('foodbag.be', 'https://www.foodbag.be/nl/home'),;
