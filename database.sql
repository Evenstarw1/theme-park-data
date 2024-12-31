-- CREATE DATABASE themepark;

-- To remove all tables
-- DROP SCHEMA public CASCADE;
-- CREATE SCHEMA public;

CREATE TABLE themeparks (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name TEXT NOT NULL,
    location POINT NOT NULL,
    description TEXT NOT NULL,
    picture TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_name UNIQUE (name)
);

CREATE TABLE categories (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE themeparks_categories (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    themepark_id INT REFERENCES themeparks (id) ON UPDATE CASCADE ON DELETE CASCADE,
    category_id INT REFERENCES categories (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE attractions (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    themepark_id INT REFERENCES themeparks (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL, -- SHA256 from the actual password
    access_level INT NOT NULL, -- 1 Admin, 2 normal user
    birth_date DATE NOT NULL,
    city TEXT NOT NULL,
    profile_picture TEXT, -- link to it?
    description TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_email UNIQUE (email)
);

CREATE TABLE users_categories (
     id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
     user_id INT REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
     category_id INT REFERENCES categories (id) ON UPDATE CASCADE ON DELETE CASCADE,
     created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE tokens (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    token TEXT NOT NULL,
    user_id INT REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE comments (
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_id INT REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    themepark_id INT REFERENCES themeparks (id) ON UPDATE CASCADE ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


-- Test dataset data

INSERT INTO themeparks (id, name, location, description, picture)
VALUES
    (1, 'Parque de atracciones de Madrid', POINT(40.412021843909045, -3.7493498990303147), 'Parque propiedad del ayuntamiento de Madrid', 'https://www.transfersandexperiences.com/images/actividades/1027/parque-de-atracciones-madrid.jpg'),
    (2, 'Parque Warner de Madrid', POINT(40.23184478750921, -3.5925680021535955), 'Parque temático con personajes de Warner Bros studios', 'https://www.transfersandexperiences.com/images/actividades/1025/parque-warner-madrid.jpeg'),
    (3, 'Port Aventura', POINT(41.09350804068493, 1.1611417895801308), 'El mejor parque de atracciones y temático de España', 'https://image.jimcdn.com/app/cms/image/transf/dimension=778x10000:format=jpg/path/seb53d1550d7485e8/image/ibcb2924f48a1dfee/version/1638261919/entradas-port-aventura.jpg');


INSERT INTO categories (id, name)
VALUES
    (1, 'Emoción'),
    (2, 'Familiar'),
    (3, 'Niños'),
    (4, 'Relax');

INSERT INTO themeparks_categories (id, themepark_id, category_id)
VALUES
    (1,1,1), -- PAM has all categories (Trilling, Family friendly, Children and Relaxing)
    (2,1,2),
    (3,1,3),
    (4,1,4),
    (5,2,1), -- Warner only 1 (thrilling) and 2 (Family friendly, Family friendly, Children and Relaxing)
    (6,2,2),
    (7,3,1); -- Port Aventura only 1 (thrilling)


INSERT INTO attractions (id, themepark_id, name)
VALUES
     (1, 1, 'La lanzadera'), -- Belongs to PAM
     (2, 1, 'El tornado'), -- Belongs to PAM
     (3, 2, 'Stunt Fall'), -- Belongs to Warner
     (4, 2, 'Batman returns'), -- Belongs to Warner
     (5, 3, 'Shambala'), -- Belongs to Port Aventura
     (6, 3, 'Dragon Khan'); -- Belongs to Port Aventura

INSERT INTO users (id, name, email, password, access_level, birth_date, city, profile_picture, description)
VALUES
    (1, 'Solid Snake', 'solid.snake@konami.jp', '178c15232b8899b70ebc1c0e9eee1de80cd5031501a9d9dc1ed31f8077f8313c', 1, '1991-10-31', 'Alaska City', 'http://www.hardcoregaming101.net/wp-content/uploads/2023/05/metal-gear-solid-3-146-1536x864.jpg', 'I love rollercoasters lol'), --password is snake1234
    (2, 'Revolver Occelot', 'revolver.occelote@konami.jp', 'ea5e443277359401ae943262d5d99c51743e2b2093b5c58340e859ce5809618a', 2, '1966-02-21', 'Moscow', 'https://image.civitai.com/xG1nkqKTMzGDvpLrqFT7WA/3af0e904-f03f-42e0-af2f-ae044fefa370/original=true,quality=90/dVykdLx0pwlQlCo7_WzFs.jpeg', 'I hate snakes'), --password is ocelote1234
    (3, 'Meryl Silverburgh', 'meryl.silverburgh@konami.jp', '85f7280ad47598661a220198e63d91836ecb5c5780eab94707e92bf299572b49', 2, '1994-09-1', 'Pariss', 'https://image.civitai.com/xG1nkqKTMzGDvpLrqFT7WA/de899f50-7d1c-4460-bfa8-2dca60b05b2a/width=450/00372-1034797549.jpeg', 'Im not a rokie!'); -- password is meryl1234


INSERT INTO users_categories (id, user_id, category_id)
VALUES
    (1,1,1), -- Solid Snake loves Trilling, Family friendly and Children theme parks
    (2,1,2),
    (3,1,3),
    (4,2,1), -- Revolver Occelot loves Trilling and Family friendly theme parks
    (5,2,2),
    (6,3,4); -- Meryl Silverburgh loves only Relaxing

INSERT INTO comments (id, user_id, themepark_id, comment)
VALUES
    (1, 1, 1, 'Me encanta la lanzadera, la adrenalina es lo más'), -- Solid Snake de PAM
    (2, 1, 1, 'Los fiordos son muy aburridos, aunque te mojas un montón'), -- Solid Snake de PAM
    (3, 2, 2, 'La nueva de batman es mi favorita'), -- Revolver Occelot de Warner
    (4, 2, 3, 'Siempre hay mucha gente'), -- Revolver Occelot de port aventura
    (5, 3, 2, 'Evitad ir en verano, las temperaturas supertan los 40 grados :S'), -- Meryl Silverburgh de Warner
    (6, 3, 3, 'Mi parque temático favorito de España'), -- Meryl Silverburgh de Port aventura
    (7, 1, 3, 'Ferrari-Land es una pasada'); -- Solid Snake de Port aventura

-- Synchronize sequences
-- N.B: We need to do this since we are setting ids manually in the test dataset. If we do not do this,
-- we would face issues with violation of primary key constraints.
-- Solution found -> https://stackoverflow.com/questions/9108833/postgres-autoincrement-not-updated-on-explicit-id-inserts
SELECT setval('attractions_id_seq', (SELECT MAX(id) from attractions));
SELECT setval('categories_id_seq', (SELECT MAX(id) from categories));
SELECT setval('comments_id_seq', (SELECT MAX(id) from comments));
SELECT setval('themeparks_categories_id_seq', (SELECT MAX(id) from themeparks_categories));
SELECT setval('themeparks_id_seq', (SELECT MAX(id) from themeparks));
SELECT setval('tokens_id_seq', (SELECT MAX(id) from tokens));
SELECT setval('users_categories_id_seq', (SELECT MAX(id) from users_categories));
SELECT setval('users_id_seq', (SELECT MAX(id) from users));
