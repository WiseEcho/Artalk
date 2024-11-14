package common

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt"
)

type userClaims struct {
	SecretID     string `json:"secret_id"`
	UserID       string `json:"user_id"`
	TokenID      string `json:"token_id"`
	TokenName    string `json:"token_name"`
	OpenID       string `json:"open_id"`
	RegisterType int8   `json:"register_type"`
}

func DecryptAuth(tokenSecret, authInfo string) (*userClaims, error) {
	if authInfo == "" {
		return nil, errors.New("auth is null")
	}

	authArray := strings.Split(authInfo, " ")
	if len(authArray) != 2 {
		return nil, errors.New("auth illegal")
	}

	// RFC 6750
	if authArray[0] != "Bearer" {
		return nil, errors.New("auth illegal")
	}

	info := struct {
		userClaims
		jwt.StandardClaims
	}{}

	token, err := jwt.ParseWithClaims(authArray[1], &info, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("param error")
	}

	claim := &userClaims{
		SecretID:     info.SecretID,
		UserID:       info.UserID,
		TokenID:      info.TokenID,
		TokenName:    info.TokenName,
		OpenID:       info.OpenID,
		RegisterType: info.RegisterType,
	}

	return claim, nil
}
