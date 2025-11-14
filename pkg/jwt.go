package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTService(accessSecret, refreshSecret []byte, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (s *JWTService) GenerateTokens(userID string) (string, string, error) {
	accessToken, err := s.generateAccessToken(userID, s.accessSecret)
	if err != nil {
		return "", "", fmt.Errorf("s.generateAccessToken: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(userID, s.refreshSecret)
	if err != nil {
		return "", "", fmt.Errorf("s.generateRefreshToken: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *JWTService) generateAccessToken(userID string, accessSecret []byte) (string, error) {
	claims := s.makeMapClaims(userID, s.accessTTL)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(accessSecret)
	if err != nil {
		return "", fmt.Errorf("jwt.NewWithClaims: %w", err)
	}

	return token, nil
}

func (s *JWTService) generateRefreshToken(userID string, refreshSecret []byte) (string, error) {
	claims := s.makeMapClaims(userID, s.refreshTTL)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(refreshSecret)
	if err != nil {
		return "", fmt.Errorf("jwt.NewWithClaims: %w", err)
	}

	return token, nil
}

func (s *JWTService) makeMapClaims(userID string, ttl time.Duration) jwt.MapClaims {
	now := time.Now()
	return jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(ttl).Unix(),
		"iat":     now.Unix(),
	}
}

func (s *JWTService) ValidateAccessToken(token string) (string, error) {
	return s.validateToken(token, s.accessSecret)
}

func (s *JWTService) ValidateRefreshToken(token string) (string, error) {
	return s.validateToken(token, s.refreshSecret)
}

func (s *JWTService) validateToken(token string, secret []byte) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("jwt.Parse: %w", err)
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", errors.New("user_id not found in token")
	}

	return "", errors.New("invalid token")
}
