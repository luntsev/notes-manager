package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	Secret             string
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
	logger             *logs.Logger
}

func NewJWT(secret string, accessExp, refreshExp int) *JWT {
	return &JWT{
		Secret:             secret,
		accessTokenExpire:  time.Duration(accessExp) * time.Second,
		refreshTokenExpire: time.Duration(refreshExp) * time.Second,
	}
}

func (j *JWT) Create(data *JWTData) (*tokens, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"exp":   time.Now().Add(j.accessTokenExpire).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"exp":   time.Now().Add(j.refreshTokenExpire).Unix(),
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
			j.logger.WriteError("Wrong jwt-token signing method")
			return nil, errors.New("wrong jwt-token signing method")
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		j.logger.WriteWarn("Invalid token")
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
