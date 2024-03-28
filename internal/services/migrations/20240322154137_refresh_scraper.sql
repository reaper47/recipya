-- +goose Up
INSERT INTO websites (host, url)
VALUES ('puurgezond.nl', 'https://www.puurgezond.nl'),
       ('jaimyskitchen.nl', 'https://jaimyskitchen.nl'),
       ('leukerecepten.nl', 'https://www.leukerecepten.nl'),
       ('bettybossi.ch', 'https://www.bettybossi.ch'),
       ('reddit.com', 'https://www.reddit.com'),
       ('marmiton.org', 'https://www.marmiton.org'),
       ('yumelise.fr', 'https://www.yumelise.fr'),
       ('lidl-kochen.de', 'https://www.lidl-kochen.de/rezeptwelt'),
       ('all-clad.com', 'https://www.all-clad.com'),
       ('francescakookt.nl', 'https://www.francescakookt.nl'),
       ('quitoque.fr', 'https://www.quitoque.fr'),
       ('bingingwithbabish.com', 'https://www.bingingwithbabish.com'),
       ('chuckycruz.substack.com', 'https://chuckycruz.substack.com'),
       ('kuchynalidla.sk', 'https://kuchynalidla.sk'),
       ('myjewishlearning.com', 'https://www.myjewishlearning.com'),
       ('drinkoteket.se', 'https://drinkoteket.se');

-- +goose Down
DELETE
FROM websites
WHERE host IN ('puurgezond.nl', 'jaimyskitchen.nl', 'leukerecepten.nl', 'bettybossi.ch', 'reddit.com',
               'marmiton.org', 'yumelise.fr', 'lidl-kochen.de', 'all-clad.com', 'francescakookt.nl', 'quitoque.fr',
               'bingingwithbabish.com', 'chuckycruz.substack.com', 'kuchynalidla.sk', 'myjewishlearning.com',
               'drinkoteket.se');
