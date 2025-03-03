package dtos

// TODO: add tags for validation
type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	APIKey  string `json:"api_key"`
}
