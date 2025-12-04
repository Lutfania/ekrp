package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	RoleID      string   `json:"role_id"`
	Token       string   `json:"token"`
	Permissions []string `json:"permissions"`
}
