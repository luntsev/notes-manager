package database

import (
	"fmt"

	"github.com/luntsev/notes-manager/auth/configs"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	Db     *gorm.DB
	logger logs.Logger
}

func NewPostgresDB(conf *configs.Config, logger *logs.Logger) *PostgresDB {
	postgresURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		conf.PostgresDataBase.PostgresHost,
		conf.PostgresDataBase.PostgresUser,
		conf.PostgresDataBase.PostgresPassword,
		conf.PostgresDataBase.PostgresDbName,
		conf.PostgresDataBase.PostgresPort,
		conf.PostgresDataBase.PostgresUseSSL)

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		logger.WriteError(fmt.Sprintf("Unable connect to database: %s", err.Error()))
		panic(err)
	}

	logger.WriteInfo("Successful connection to the database")
	return &PostgresDB{
		Db:     db,
		logger: *logger,
	}
}
