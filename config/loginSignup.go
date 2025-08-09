package config

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
