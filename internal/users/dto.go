package users

type CreateUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
