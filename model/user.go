package model

// @Description Пользователь - существует постоянно
type User struct {
	Id        int    `json:"-" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" binding:"required" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" binding:"required" db:"password_hash"`
}
