package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"notes-manager/configs"
	"notes-manager/envs"
	"notes-manager/pkg/database"
	"notes-manager/pkg/logs"
	"notes-manager/repository"
)

func InitServer() *gin.Engine {
	envs.LoadEnvs()

	conf, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := logs.NewLogger(conf)
	notesDb := database.NewMongoDB(conf, logger)
	cache, err := database.NewRedisDB(conf, logger)
	if err != nil {
		logger.WriteError(fmt.Sprintf("Error when cache init: %s", err.Error()))
		panic(err.Error())
	}

	notesRepo := repository.NewNotesRepository(conf, notesDb, cache, logger)

	return InitRoutes(notesRepo, logger)
}

func StartServer(router *gin.Engine) {
	router.Run(":9100")
}
