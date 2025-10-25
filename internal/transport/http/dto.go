package http

type UserDTO struct {
	ID string `json:"id"`
	UserName string `json:"user_name"`
	Email string  `json:"email"`
}

type RegisterRequest struct {
	UserName string `json:"user_name"`
	Email string  `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User UserDTO
}

type LoginRequest struct {
	Email string  `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User UserDTO
}

type UpdateRequest struct {
	UserName string `json:"user_name"`
	Email string  `json:"email"`
}

type UpdateResponse struct {
	User UserDTO
}