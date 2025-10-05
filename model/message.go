package model

import "time"

// @Description Сообшения чата
type Message struct {
	Id        int       `json:"id"`
	Room      int       `json:"room" db:"room_id"`
	User      int       `json:"user" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (r *Room) GetIdMes() int {
	return r.Id
}

func (u *User) GetIdUsMes() int {
	return u.Id
}
