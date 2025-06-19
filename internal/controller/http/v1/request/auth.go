package request

type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"admin"`
	Password string `json:"password" validate:"required" example:"admin123"`
}
