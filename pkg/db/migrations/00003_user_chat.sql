-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_chat (
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    chat_id INT REFERENCES chats(id) ON DELETE CASCADE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_chat;
-- +goose StatementEnd
