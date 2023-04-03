CREATE TABLE IF NOT EXISTS nix_chat.rooms (
    id              serial PRIMARY KEY,
    name         varchar(100) NOT NULL,
    private TINYINT NULL
);
