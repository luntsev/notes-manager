package repository

import (
	"fmt"

	"github.com/luntsev/notes-manager/auth/models"
	"github.com/luntsev/notes-manager/auth/pkg/database"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"gorm.io/gorm/clause"
)

type AuthRepository struct {
	DataBase *database.PostgresDB
	logger   *logs.Logger
}

func NewAuthRepository(db *database.PostgresDB, logger *logs.Logger) *AuthRepository {
	return &AuthRepository{
		DataBase: db,
		logger:   logger,
	}
}

func (repo *AuthRepository) Create(user *models.User) error {
	result := repo.DataBase.Create(user)
	if result.Error != nil {
		repo.logger.WriteError(fmt.Sprintf("Unable create user record: %s", result.Error.Error()))
	} else {
		repo.logger.WriteInfo("User record is created")
	}
	return result.Error
}

func (repo *AuthRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	result := repo.DataBase.First(&user, "email = ?", email)
	if result.Error != nil {
		repo.logger.WriteWarn(fmt.Sprintf("Unable read user record in database: %s", result.Error.Error()))
		return nil, result.Error
	}

	return &user, nil
}

func (repo *AuthRepository) Update(user *models.User) error {
	result := repo.DataBase.Clauses(clause.Returning{}).Updates(user)
	if result.Error != nil {
		repo.logger.WriteWarn(fmt.Sprintf("Unable update user record in database: %s", result.Error.Error()))
	}

	return result.Error
}
