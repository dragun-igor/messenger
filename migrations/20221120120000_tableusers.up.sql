CREATE TABLE IF NOT EXISTS users (
    login text NOT NULL UNIQUE,
    name text NOT NULL UNIQUE,
    password text NOT NULL
);
