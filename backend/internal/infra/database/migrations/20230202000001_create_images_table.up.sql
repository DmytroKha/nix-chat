CREATE TABLE IF NOT EXISTS nix_chat.images
(
    id           serial PRIMARY KEY,
    user_id      int          NOT NULL,
    name         varchar(250) NOT NULL,
    created_date timestamp    NOT NULL,
    updated_date timestamp    NOT NULL,
    deleted bool    NOT NULL
);
