-- +goose Up
INSERT INTO websites (host, url) 
VALUES ('gutekueche.at', 'https://www.gutekueche.at/');

-- +goose Down
DELETE FROM websites
WHERE host IN ('gutekueche.at');