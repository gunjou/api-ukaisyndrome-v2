package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Platform string `json:"platform"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SessionEntity struct {
	UserID   int
	Role     string
	Platform string
	Token    string
}

type UserAuthEntity struct {
	ID             int
	Email          string
	Password       string
	KodePemulihan  string
	Role           string
	Status         int
}