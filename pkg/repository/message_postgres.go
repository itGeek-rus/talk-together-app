package repository

import (
	"fmt"
	"github.com/firstproject/talk-together-app/model"
	"github.com/jmoiron/sqlx"
)

type MessagePostgres struct {
	db *sqlx.DB
}

func NewMessagePostgres(db *sqlx.DB) *MessagePostgres {
	return &MessagePostgres{db: db}
}

func (r *MessagePostgres) CreateMessage(roomId, userId int, content string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (room_id, user_id, content) VALUES ($1, $2, $3) RETURNING id", messagesTable)
	err := r.db.Get(&id, query, roomId, userId, content)

	return id, err
}

func (r *MessagePostgres) GetRoomMessages(roomId int) ([]model.Message, error) {
	var messages []model.Message
	query := fmt.Sprintf(`
						SELECT m.id, m.room_id, m.user_id, m.content, m.created_at
						FROM %s m WHERE m.room_id = $1
						ORDER BY m.created_at`, messagesTable)

	err := r.db.Select(&messages, query, roomId)
	return messages, err
}

func (r *MessagePostgres) DeleteMessage(messageId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", messagesTable)

	result, err := r.db.Exec(query, messageId, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

func (r *MessagePostgres) GetMessageOwener(messageId int) (int, error) {
	var userId int

	query := fmt.Sprintf("SELECT user_id FROM %s WHERE id = $1", messagesTable)

	err := r.db.Get(&userId, query, messageId)
	return userId, err
}

func (r *MessagePostgres) UpdateMessage(messageId, userId int, content string) error {
	query := fmt.Sprintf("UPDATE %s SET content = $1 WHERE id = $2 AND user_id = $3", messagesTable)
	result, err := r.db.Exec(query, content, messageId, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

func (r *MessagePostgres) GetMessageById(messageId int) (model.Message, error) {
	var message model.Message

	query := fmt.Sprintf(`
			SELECT m.id, m.room_id, m.user_id, m.content, m.created_at
			FROM %s m WHERE m.id = $1`, messagesTable)

	err := r.db.Get(&message, query, messageId)

	return message, err
}
