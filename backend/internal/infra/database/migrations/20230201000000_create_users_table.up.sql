CREATE TABLE IF NOT EXISTS nix_chat.users (
    id              serial PRIMARY KEY,
    password        varchar(100) NOT NULL,
    name         varchar(100) NOT NULL
);
