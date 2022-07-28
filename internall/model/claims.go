package model

import "github.com/golang-jwt/jwt"

type AccessClaims struct {
	UserId   int    `json:"user-id"`
	UserName string `json:"user-name"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	UserId int `json:"user-id"`
	jwt.StandardClaims
}
