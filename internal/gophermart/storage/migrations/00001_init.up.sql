BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS users
(
    id serial primary key,
    login varchar(255) not null unique,
    password_hash varchar(255) not null
);
COMMIT;
