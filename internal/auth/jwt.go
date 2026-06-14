package auth

import (
	"time"

	"go-chat/internal/config"

	"github.com/golang-jwt/jwt/v5"
)
type AuthService struct {
      secret      []byte
      expiryHours int
}

func NewAuthService(auth config.AuthConfig) *AuthService {
	return &AuthService{secret: []byte(auth.JWT_SECRET)}
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.JwtExpiryHours) * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secret)
}