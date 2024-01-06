-- +goose Up
DELETE
FROM websites
WHERE host = 'mob.co.uk'
   OR host = 'nutritionbynathalie.com';

-- +goose Down
INSERT INTO websites (host, url)
VALUES ('mob.co.uk', 'https://www.mob.co.uk'),
       ('nutritionbynathalie.com', 'https://www.nutritionbynathalie.com');
