-- +goose Up
DELETE
FROM websites
WHERE host = 'mob.co.uk';

-- +goose Down
INSERT INTO websites (host, url)
VALUES ('mob.co.uk', 'https://www.mob.co.uk');
