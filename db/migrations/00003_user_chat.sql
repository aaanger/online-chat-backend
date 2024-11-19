-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_room (
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    room_id INT REFERENCES room(id) ON DELETE CASCADE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_room;
-- +goose StatementEnd
