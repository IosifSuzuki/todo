package model

import "time"

type AccountModel struct {
	Id           int       `json:"id"`
	UserName     string    `json:"user-name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created-at"`
	HashPassword string    `json:"-"`
}
