package repository

import (
	"fmt"
	"github.com/firstproject/talk-together-app/model"
	"github.com/jmoiron/sqlx"
)

type ClientPostgres struct {
	db *sqlx.DB
}

func NewClientPostgres(db *sqlx.DB) *ClientPostgres {
	return &ClientPostgres{db: db}
}

func (r *ClientPostgres) AddClientToRoom(roomId, userId int) error {
	query := fmt.Sprintf("INSERT INTO %s (room_id, user_id) VALUES ($1, $2)", clientsTable)
	_, err := r.db.Exec(query, roomId, userId)
	return err
}

func (r *ClientPostgres) RemoveClientFromRoom(roomId, userId int) error {
	query := fmt.Sprintf("UPDATE %s SET disconnected_at = NOW() WHEERE room_id = $1 AND user_id = $2 AND disconnected_at IS NULL", clientsTable)
	_, err := r.db.Exec(query, roomId, userId)
	return err
}

func (r *ClientPostgres) GetRoomClients(roomId int) ([]model.User, error) {
	var users []model.User

	query := fmt.Sprintf(`SELECT u.* FROM %s u
						INNER JOIN %s c ON u.id = c.user_id
						WHERE c.room_id = $1 AND c.disconnected_at IS NULL`, usersTable, clientsTable)
	err := r.db.Select(&users, query, roomId)
	return users, err
}
