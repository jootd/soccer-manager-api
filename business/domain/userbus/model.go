package userbus

import "time"

type User struct {
	Username    string
	Password    string
	TeamID      int
	DateCreated time.Time
	DateUpdated time.Time
}

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TeamID   int    `json:"-"`
}

type UpdateUser struct {
	Username     *string
	PasswordHash *string
	TeamID       *int
}
