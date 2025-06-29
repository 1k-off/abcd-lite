package domain

type LoginRequest struct {
	AdminToken string `json:"admin_token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
