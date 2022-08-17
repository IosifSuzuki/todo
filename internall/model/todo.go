package model

import "time"

type Todo struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"created-on"`
	UpdatedOn   time.Time `json:"updated-on"`
	Closed      bool      `json:"closed"`
}
