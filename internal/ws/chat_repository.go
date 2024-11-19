package ws

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	apperrors "onlineChat/internal/errors"
)

type ChatRepository struct {
	db *sqlx.DB
}

func (r *ChatRepository) CreateChat(name string, userID int) (*Chat, error) {
	chat := Chat{
		Name: name,
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repository create room: %w", err)
	}

	row := tx.QueryRow(`INSERT INTO chats (name) VALUES($1) RETURNING id;`, name)
	err = row.Scan(&chat.ID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("repository create room: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO user_chat (user_id, chat_id) VALUES($1, $2);`, userID, chat.ID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("repository create room: %w", err)
	}
	return &chat, tx.Commit()
}

func (r *ChatRepository) SaveChat(userID int, message *Message) error {
	_, err := r.db.Exec(`INSERT INTO chats (user_id, message) VALUES($1, $2)`, userID, message.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepository) JoinChat(userID, chatID int) (*Chat, error) {
	var chat Chat

	query := fmt.Sprintf(`SELECT * FROM chats c INNER JOIN user_chat uc ON c.id = uc.room_id WHERE uc.user_id=$1 AND uc.room_id=$2;`)
	err := r.db.Get(&chat, query, userID, chatID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ChatNotFound
		}
		return nil, err
	}

	return &chat, err
}
