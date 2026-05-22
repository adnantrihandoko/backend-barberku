package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AdminUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	PINHash   string    `json:"-"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email string `json:"email"`
	PIN   string `json:"pin"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  *AdminUser `json:"user"`
}

type JWTClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Role   UserRole `json:"role"`
	jwt.RegisteredClaims
}
