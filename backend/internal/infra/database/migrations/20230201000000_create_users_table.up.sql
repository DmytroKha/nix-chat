CREATE TABLE IF NOT EXISTS nix_chat.users (
    id              serial PRIMARY KEY,
    uid              varchar(100) NOT NULL,
    password        varchar(100) NOT NULL,
    name         varchar(100) NOT NULL
);
