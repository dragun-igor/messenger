CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    login text NOT NULL UNIQUE,
    name text NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id integer PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    time timestamp NOT NULL,
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

CREATE TABLE IF NOT EXISTS friends (
    user1 text NOT NULL,
    user2 text NOT NULL,
    UNIQUE (user1, user2),
    CONSTRAINT fk_user1
        FOREIGN KEY(user1)
            REFERENCES users(name),
    CONSTRAINT fk_user2
        FOREIGN KEY(user2)
            REFERENCES users(name)
);

CREATE TABLE IF NOT EXISTS requests_to_friends_list (
    requester text NOT NULL,
    receiver text NOT NULL,
    UNIQUE (requester, receiver),
    CONSTRAINT fk_requester
        FOREIGN KEY(requester)
            REFERENCES users(name),
    CONSTRAINT fk_receiver
        FOREIGN KEY(receiver)
            REFERENCES users(name)
);