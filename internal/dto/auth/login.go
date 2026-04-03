package auth

// RegisterRequest is the DTO for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest is the DTO for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is returned after successful login or token refresh
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshRequest is the DTO for refreshing tokens
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
