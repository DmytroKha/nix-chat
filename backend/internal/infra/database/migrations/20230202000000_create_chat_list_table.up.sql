CREATE TABLE IF NOT EXISTS nix_chat.chat_list (
    id              serial PRIMARY KEY,
    name         varchar(100) NOT NULL,
    user_id           int NOT NULL,
    client_id           int NOT NULL,
    created_date timestamp    NOT NULL,
    updated_date timestamp    NOT NULL,
    deleted_date timestamp    NULL
);
