package model

import (
	"errors"
	"strings"
	"time"
)

// @Description Комнаты чата
type Room struct {
	Id          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UpdateRoomInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (i UpdateRoomInput) Validate() error {
	if i.Name == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	if i.Name != nil && strings.TrimSpace(*i.Name) == " " {
		return errors.New("name cannot be empty")
	}

	if i.Description != nil && strings.TrimSpace(*i.Description) == " " {
		return errors.New("description cannot be empty")
	}

	return nil
}

func (u *User) GetId() int {
	return u.Id
}
