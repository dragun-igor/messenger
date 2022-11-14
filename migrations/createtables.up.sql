CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    first_name text,
    second_name text,
    login_name text NOT NULL UNIQUE,
    pswd text NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id integer PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    sender_id integer NOT NULL,
    receiver_id integer NOT NULL,
    msg text NOT NULL,
    CONSTRAINT fk_sender_id
        FOREIGN KEY(sender_id)
            REFERENCES users(id),
    CONSTRAINT fk_receiver_id
        FOREIGN KEY(receiver_id)
            REFERENCES users(id)
);