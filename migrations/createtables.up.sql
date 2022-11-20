CREATE TABLE IF NOT EXISTS users (
    login text NOT NULL UNIQUE,
    name text NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    sender text NOT NULL,
    receiver text NOT NULL,
    msg text NOT NULL,
    CONSTRAINT fk_sender
        FOREIGN KEY(sender)
            REFERENCES users(name),
    CONSTRAINT fk_receiver
        FOREIGN KEY(receiver)
            REFERENCES users(name)
);