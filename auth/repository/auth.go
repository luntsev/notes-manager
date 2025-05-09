package repository

import (
	"fmt"

	"github.com/luntsev/notes-manager/auth/models"
	"github.com/luntsev/notes-manager/auth/pkg/database"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
)

type AuthRepositury struct {
	DataBase *database.PostgresDB
	logger   *logs.Logger
}

func NewAuthRepository(db *database.PostgresDB, logger *logs.Logger) *AuthRepositury {
	return &AuthRepositury{
		DataBase: db,
		logger:   logger,
	}
}

func (repo *AuthRepositury) Create(user *models.User) error {
	result := repo.DataBase.Create(user)
	if result.Error != nil {
		repo.logger.WriteError(fmt.Sprintf("Unable create user record: %s", result.Error.Error()))
	} else {
		repo.logger.WriteInfo("User record is created")
	}
	return result.Error
}
