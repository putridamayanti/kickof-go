package models

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type UserRequest struct {
	Active     bool        `json:"active"`
	Email      string      `json:"email"`
	Image      string      `json:"image"`
	Main       bool        `json:"main"`
	Name       string      `json:"name"`
	Owner      bool        `json:"owner"`
	Password   string      `json:"password"`
	Permission interface{} `json:"permission"`
	Status     bool        `json:"status"`
}
