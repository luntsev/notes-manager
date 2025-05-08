package server

import (
	"fmt"
	"log"
	"notes-manager/configs"
	"notes-manager/pkg/database"
	"notes-manager/pkg/logs"
	"notes-manager/repository"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
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

	noteRouter := NewRouter(notesRepo, logger)
	return &Server{router: noteRouter.router}
}

func (s *Server) Start(port int) {
	if port <= 0 {
		log.Fatalf("Invalid port: %d", port)
	}
	s.router.Run(fmt.Sprintf(":%d", port))
}
