package domain

type DefaultErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type DefaultInfoResponse struct {
	Message string `json:"message"`
}
