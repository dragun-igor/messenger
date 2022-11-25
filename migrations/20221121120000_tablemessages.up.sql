CREATE TABLE IF NOT EXISTS messages (
    sender varchar(32) NOT NULL,
    receiver varchar(32) NOT NULL,
    msg varchar(256) NOT NULL,
    CONSTRAINT fk_sender
        FOREIGN KEY(sender)
            REFERENCES users(name),
    CONSTRAINT fk_receiver
        FOREIGN KEY(receiver)
            REFERENCES users(name)
);