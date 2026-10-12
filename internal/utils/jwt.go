package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

// Claims adalah payload JWT aplikasi.
type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager menandatangani dan memverifikasi access & refresh token (HS256).
type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewTokenManager(accessSecret, refreshSecret string) *TokenManager {
	return &TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

// RefreshTTL dipakai untuk menentukan umur cookie refresh token.
func (m *TokenManager) RefreshTTL() time.Duration { return refreshTokenTTL }

func (m *TokenManager) GenerateAccessToken(userID uint, role string) (string, error) {
	return m.sign(m.accessSecret, userID, role, accessTokenTTL)
}

func (m *TokenManager) GenerateRefreshToken(userID uint, role string) (string, error) {
	return m.sign(m.refreshSecret, userID, role, refreshTokenTTL)
}

func (m *TokenManager) ParseAccessToken(token string) (*Claims, error) {
	return m.parse(token, m.accessSecret)
}

func (m *TokenManager) ParseRefreshToken(token string) (*Claims, error) {
	return m.parse(token, m.refreshSecret)
}

func (m *TokenManager) sign(secret []byte, userID uint, role string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func (m *TokenManager) parse(tokenString string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) { return secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
