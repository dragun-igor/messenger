CREATE TABLE IF NOT EXISTS users (
    login varchar(32) NOT NULL UNIQUE,
    name varchar(32) NOT NULL UNIQUE,
    password varchar(128) NOT NULL
);
