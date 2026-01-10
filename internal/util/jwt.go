package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secret           []byte
	accessExpiration time.Duration
	refreshExpiration time.Duration
	issuer           string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secret string, accessExpiration, refreshExpiration time.Duration) *JWTManager {
	return &JWTManager{
		secret:           []byte(secret),
		accessExpiration: accessExpiration,
		refreshExpiration: refreshExpiration,
		issuer:           "knowledge-garden",
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token for a user
func (j *JWTManager) GenerateAccessToken(userID, email string) (string, error) {
	return j.generateToken(userID, email, "access", j.accessExpiration)
}

// GenerateRefreshToken generates a refresh token for a user
func (j *JWTManager) GenerateRefreshToken(userID, email string) (string, error) {
	return j.generateToken(userID, email, "refresh", j.refreshExpiration)
}

// generateToken generates a JWT token
func (j *JWTManager) generateToken(userID, email, tokenType string, expiration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(expiration)

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			Audience:  []string{"knowledge-garden-api"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, fmt.Errorf("expected access token, got %s", claims.TokenType)
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("expected refresh token, got %s", claims.TokenType)
	}

	return claims, nil
}

// GetAccessExpiration returns the access token expiration in seconds
func (j *JWTManager) GetAccessExpiration() int64 {
	return int64(j.accessExpiration.Seconds())
}
