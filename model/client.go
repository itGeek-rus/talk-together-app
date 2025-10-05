package model

import "github.com/gorilla/websocket"

// @Description Клиент - существует во время соединения
type Client struct {
	Id   int             `json:"id"`
	Conn *websocket.Conn `json:"conn"`
	Room int             `json:"room_id" db:"room_id"`
	User int             `json:"user" db:"user_id"`
	Send chan []byte     `json:"send"`
}

func (r *Room) GetId() int {
	return r.Id
}

func (u *User) GetIdClient() int {
	return u.Id
}
