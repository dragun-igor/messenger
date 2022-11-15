CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    login text NOT NULL UNIQUE,
    name text NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id integer PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY,
    time timestamp NOT NULL,
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

CREATE TABLE IF NOT EXISTS friends (
    user1 integer NOT NULL,
    user2 integer NOT NULL,
    CONSTRAINT fk_user1
        FOREIGN KEY(user1)
            REFERENCES users(id),
    CONSTRAINT fk_user2
        FOREIGN KEY(user2)
            REFERENCES users(id)
)

CREATE TABLE IF NOT EXISTS requests_to_friends_list (
    requester integer NOT NULL,
    receiver integer NOT NULL,
    CONSTRAINT fk_requester
        FOREIGN KEY(requester)
            REFERENCES users(id),
    CONSTRAINT fk_receiver
        FOREIGN KEY(receiver)
            REFERENCES users(id)    
)