package auth

import (
	"errors"
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

func (s *AuthService) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}
