package models

// LogoutRequest maps to logout api call params
type LogoutRequest struct {
	Email string `json:"email" validate:"required,email"`
}
