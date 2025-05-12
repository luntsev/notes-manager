package models

import (
	"errors"
	"strings"
	"time"
)

type birthdayDate struct {
	time.Time
}

func (t *birthdayDate) UnmarshalJSON(b []byte) error {
	layout := "02.01.2006"
	dateStr := strings.Trim(string(b), "\"") // remove quotes
	if dateStr == "null" {
		return errors.New("date is empty")
	}
	parsTime, err := time.Parse(layout, dateStr)
	t.Time = parsTime
	return err
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	LoginRequest
	Name     string       `json:"name" binding:"required"`
	BirthDay birthdayDate `json:"birthDay,format:2006.01.02" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateRequest struct {
	Name     *string       `json:"name,omitempty"`
	Password *string       `json:"password,omitempty"`
	BirthDay *birthdayDate `json:"birthDay,omitempty"`
}
