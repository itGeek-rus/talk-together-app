package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/firstproject/talk-together-app/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type RoomPostgres struct {
	db *sqlx.DB
}

func NewRoomPostgres(db *sqlx.DB) *RoomPostgres {
	return &RoomPostgres{db: db}
}

func (r *RoomPostgres) CreateRoom(userId int, room model.Room) (int, error) {
	var userExists bool
	checkUserQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", usersTable)
	err := r.db.Get(&userExists, checkUserQuery, userId)
	if err != nil {
		return 0, err
	}

	if !userExists {
		return 0, fmt.Errorf("user with id %d does not exist", userId)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createRoomQuery := fmt.Sprintf("INSERT INTO %s (name, description, created_by) VALUES ($1, $2, $3) RETURNING id", roomsTable)
	row := tx.QueryRow(createRoomQuery, room.Name, room.Description, userId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()

}

func (r *RoomPostgres) GetAllRooms(userId int) ([]model.Room, error) {
	var rooms []model.Room

	query := fmt.Sprintf("SELECT * FROM %s WHERE created_by = $1", roomsTable)
	err := r.db.Select(&rooms, query, userId)

	return rooms, err
}

func (r *RoomPostgres) SearchRoomByName(name string) ([]model.Room, error) {
	var rooms []model.Room

	query := fmt.Sprintf(`SELECT * FROM %s WHERE name ILIKE $1 ORDER BY name`, roomsTable)

	searchPattern := "%" + name + "%"

	err := r.db.Select(&rooms, query, searchPattern)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RoomPostgres) GetRoomById(id int) (model.Room, error) {
	var room model.Room

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", roomsTable)

	err := r.db.Get(&room, query, id)

	return room, err
}

func (r *RoomPostgres) UpdateRoom(roomId, userId int, input model.UpdateRoomInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if len(setValues) == 0 {
		return errors.New("no fields to update")
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d AND created_by = $%d", roomsTable, setQuery, argId, argId+1)

	args = append(args, roomId, userId)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RoomPostgres) DeleteRoom(userId, roomId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND created_by = $2", roomsTable)

	result, err := r.db.Exec(query, roomId, userId)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}
