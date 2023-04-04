CREATE TABLE IF NOT EXISTS nix_chat.friend_list (
    id              serial PRIMARY KEY,
    user_id           int NOT NULL,
    friend_id           int NOT NULL,
    room_id           int NOT NULL
);
