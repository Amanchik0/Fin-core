package models

type AuthService interface {
	ValidateToken(token string) (*AuthUser, error)
	GetUserByID(UserID string) (*AuthUser, error)
}
