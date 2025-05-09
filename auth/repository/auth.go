package repository

import (
	"fmt"
	"github.com/luntsev/notes-manager/auth/models"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"gorm.io/gorm"
)

type AuthRepositury struct {
	db     *gorm.DB
	logger *logs.Logger
}

func NewAuthRepository(db *gorm.DB, logger logs.Logger) *AuthRepositury {
	return &AuthRepositury{
		db:     &gorm.DB{},
		logger: &logger,
	}
}

func (repo *AuthRepositury) Create(user *models.User) error {
	result := repo.db.Create(user)
	if result.Error != nil {
		repo.logger.WriteError(fmt.Sprintf("Unable create user record: %s", result.Error.Error()))
	} else {
		repo.logger.WriteInfo("User record is created")
	}
	return result.Error
}
