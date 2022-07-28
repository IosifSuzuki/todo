package request

type RegistrationForm struct {
	UserName string `json:"user-name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegistrationForm) IsValidated() bool {
	return true
}
