-- +goose Up
INSERT INTO websites (host, url)
VALUES ('recipecommunity.com.au', 'https://www.recipecommunity.com.au'),
       ('livingthegreenlife.com', 'https://livingthegreenlife.com'),
       ('aberlehome.com', 'https://aberlehome.com'),
       ('argiro.gr', 'https://www.argiro.gr');

-- +goose Down
DELETE
FROM websites
WHERE host IN
      ('recipecommunity.com.au', 'livingthegreenlife.com', 'aberlehome.com', 'argiro.gr');
