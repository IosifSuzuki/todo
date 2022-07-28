package model

type Ping struct {
	UserId  int    `json:"user-id"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
