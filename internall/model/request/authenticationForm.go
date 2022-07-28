package request

type AuthenticationForm struct {
	UserName string `json:"user-name"`
	Password string `json:"password"`
}
