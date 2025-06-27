package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/WBianchi/maiscrianca/configs"
	"github.com/WBianchi/maiscrianca/models"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims representa os claims do token JWT
type JWTClaims struct {
	UserID string      `json:"userId"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken gera um novo token JWT para o usuário
func GenerateToken(user *models.User, config *configs.Config) (string, error) {
	expirationTime := time.Now().Add(time.Hour * config.JWTExpirationHours)
	
	claims := &JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateToken valida um token JWT e retorna os claims
func ValidateToken(tokenString string, config *configs.Config) (*JWTClaims, error) {
	claims := &JWTClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, errors.New("token inválido")
	}
	
	return claims, nil
}
