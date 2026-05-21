package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("email atau PIN salah")
	ErrAccountInactive    = errors.New("akun tidak aktif")
	ErrInvalidPINFormat   = errors.New("PIN harus 4-6 digit angka")
)

type AuthServiceImpl struct {
	authRepo repository.AuthRepository
	jwtSecret []byte
}

func NewAuthService(authRepo repository.AuthRepository, jwtSecret string) *AuthServiceImpl {
	return &AuthServiceImpl{
		authRepo:  authRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, pin string) (*entity.LoginResponse, error) {
	if !isValidPIN(pin) {
		return nil, ErrInvalidPINFormat
	}

	admin, err := s.authRepo.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !admin.IsActive {
		return nil, ErrAccountInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PINHash), []byte(pin)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(admin)
	if err != nil {
		return nil, err
	}

	return &entity.LoginResponse{
		Token: token,
		User:  admin,
	}, nil
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*entity.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(*entity.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}

func (s *AuthServiceImpl) generateToken(admin *entity.AdminUser) (string, error) {
	claims := entity.JWTClaims{
		UserID: admin.ID,
		Email:  admin.Email,
		Role:   admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func isValidPIN(pin string) bool {
	if len(pin) < 4 || len(pin) > 6 {
		return false
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func HashPIN(pin string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
