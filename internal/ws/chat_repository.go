package ws

import (
	"database/sql"
	"errors"
	"fmt"
)

type IChatRepository interface {
	CreateChat(name string, userID int) (*Chat, error)
	SaveChat(chatID, userID int, message *Message) error
	JoinChat(userID, chatID int) (*Chat, error)
}

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) CreateChat(name string, userID int) (*Chat, error) {
	chat := Chat{
		Name: name,
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("repository create room: %w", err)
	}

	row := tx.QueryRow(`INSERT INTO chats (name, created_at) VALUES($1, current_timestamp) RETURNING id;`, name)
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

func (r *ChatRepository) SaveChat(chatID, userID int, message *Message) error {
	_, err := r.db.Exec(`INSERT INTO messages (chat_id, user_id, content, created_at) VALUES($1, $2, $3, current_timestamp)`, chatID, userID, message.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepository) JoinChat(userID, chatID int) (*Chat, error) {
	var chat Chat

	row := r.db.QueryRow(`SELECT * FROM chats c INNER JOIN user_chat uc ON c.id = uc.chat_id WHERE uc.user_id=$1 AND uc.chat_id=$2;`, userID, chatID)
	err := row.Scan(&chat.ID, &chat.Name, &chat.Message, &chat.CreatedAt, &chat.Clients)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("the chat could not be found")
		}
		return nil, fmt.Errorf("join chat: %w", err)
	}

	return &chat, nil
}
