package dto

import "github.com/toanuitt/bookmark_service/internal/model"

// ===== Register =====
type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,gte=8"`
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

type RegisterUserData struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	UpdatedAt   string `json:"updated_at"`
}

type RegisterResponse struct {
	Data    *RegisterUserData `json:"data"`
	Message string            `json:"message"`
}

// ===== Login =====
// LoginInputRequest represents the request payload for user login
type LoginInputRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=8"`
}

// LoginResponse represents the response body for user login
type LoginResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

// ===== Get self Info =====
// ProfileResponse represents the response body for user profile
type ProfileResponse struct {
	Data *model.User `json:"data"`
}

// ===== Update self Info =====
// UpdateProfileRequest represents the request payload for profile updates
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,email"`
}
