package model

type Credentials struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
}
