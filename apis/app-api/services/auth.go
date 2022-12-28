package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/loadfms/serverless-golang-monorepo/apis/app-api/types"
)

type Auth struct {
	tokenSecret []byte
}

func NewAuthService() (*Auth, error) {
	return &Auth{
		tokenSecret: []byte("super-secret-key-here"),
	}, nil
}

func (s *Auth) GenerateJWT(userPK string, durationHour int, tokenType string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Duration(durationHour) * time.Hour).Unix()
	claims["authorized"] = true
	claims["userPK"] = userPK
	claims["tokenType"] = tokenType

	tokenString, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Auth) ValidateJWT(receivedToken string, expectedTokenType string) (context types.AuthContext, err error) {
	token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
		}
		return []byte(s.tokenSecret), nil
	})

	if err != nil {
		return context, err
	}

	if !token.Valid {
		return context, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return context, err
	}

	userPK := claims["userPK"].(string)
	tokenType := claims["tokenType"].(string)

	if tokenType != expectedTokenType {
		return context, errors.New("Sessão expirada")
	}

	return types.AuthContext{
		PK:        userPK,
		TokenType: tokenType,
	}, nil
}
