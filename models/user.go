package models

import "time"

type User struct {
	Name     *string    `json:"name"`
	Username string     `json:"username"`
	Email    *string    `json:"email"`
	Password *string    `json:"password"`
	Created  *time.Time `json:"created"`
}
