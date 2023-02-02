CREATE TABLE IF NOT EXISTS nix_chat.message_list
(
    id           serial PRIMARY KEY,
    parrent_message_id      int          NOT NULL,
    chat_id      int          NOT NULL,
    user_id      int          NOT NULL,
    message      text NOT NULL,
    created_date timestamp    NOT NULL,
    updated_date timestamp    NOT NULL,
    deleted_date timestamp    NULL
);
