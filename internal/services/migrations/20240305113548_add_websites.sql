-- +goose Up
INSERT INTO websites (host, url)
VALUES ('recipecommunity.com.au', 'https://www.recipecommunity.com.au'),
       ('livingthegreenlife.com', 'https://livingthegreenlife.com'),
       ('aberlehome.com', 'https://aberlehome.com'),
       ('argiro.gr', 'https://www.argiro.gr'),
       ('bergamot.app', 'https://www.bergamot.app'),
       ('moulinex.fr', 'https://www.moulinex.fr/recette'),
       ('foodbag.be', 'https://www.foodbag.be/nl/home'),
       ('15gram.be', 'https://15gram.be/recepten'),
       ('thecookingguy.com', 'https://www.thecookingguy.com'),
       ('brianlagerstrom.com', 'https://www.brianlagerstrom.com');

-- +goose Down
DELETE
FROM websites
WHERE host IN
      ('recipecommunity.com.au', 'livingthegreenlife.com', 'aberlehome.com', 'argiro.gr', 'bergamot.app',
       'moulinex.fr', 'foodbag.be', '15gram.be', 'thecookingguy.com', 'brianlagerstrom.com');
