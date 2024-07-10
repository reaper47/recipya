-- +goose Up
DELETE
FROM websites
WHERE host IN (
    'bingingwithbabish.com', -- Removed because the website now requires membership to view recipes.
    'chuckycruz.substack.com' -- Removed because the website now requires membership to view recipes.
    );


-- +goose Down
INSERT INTO websites (host, url)
VALUES ('bingingwithbabish.com', 'https://www.bingingwithbabish.com'),
       ('chuckycruz.substack.com', 'https://chuckycruz.substack.com'),;
