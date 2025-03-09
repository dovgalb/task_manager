package user

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	UserID int    `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}
