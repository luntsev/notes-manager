package jwt

import (
	"errors"
	"time"

	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
)

type JWTData struct {
	Email string
}

type tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type JWT struct {
	Secret string
	logger *logs.Logger
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data *JWTData) (*tokens, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	signedAccessToken, err := accessToken.SignedString([]byte(j.Secret))
	if err != nil {
		j.logger.WriteError(fmt.Sprintf("Unable create JWT access token: %s", err.Error()))
		return nil, err
	}

	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		j.logger.WriteError(fmt.Sprintf("Unable create JWT refresh token: %s", err.Error()))
		return nil, err
	}

	tokens := &tokens{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokens, nil
}

func (j *JWT) Verify(tokenStr string) (*JWTData, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			j.logger.WriteError("Wrong token signing method")
			return nil, errors.New("wrong token signing method")
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		j.logger.WriteWarn("Шnvalid token")
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		j.logger.WriteError("No claims in token")
		return nil, errors.New("no claims in token")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		j.logger.WriteError("No email in token")
		return nil, errors.New("no email in token")
	}

	return &JWTData{Email: email}, nil
}
