CREATE TABLE IF NOT EXISTS nix_chat.black_list (
    id              serial PRIMARY KEY,
    user_id           int NOT NULL,
    foe_id           int NOT NULL
);
