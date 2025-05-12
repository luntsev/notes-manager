package server

import (
	"fmt"

	"github.com/luntsev/notes-manager/auth/pkg/jwt"
	"github.com/luntsev/notes-manager/notes/configs"
	"github.com/luntsev/notes-manager/notes/pkg/database"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"github.com/luntsev/notes-manager/notes/repository"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port   int
}

func NewServer() *Server {
	conf, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := logs.NewLogger(conf.LogLevel)
	notesDb := database.NewMongoDB(conf, logger)
	cache, err := database.NewRedisDB(conf, logger)
	if err != nil {
		logger.WriteError(fmt.Sprintf("Error when cache init: %s", err.Error()))
		panic(err.Error())
	}

	notesRepo := repository.NewNotesRepository(conf, notesDb, cache, logger)

	jwtServ := jwt.NewJWT(conf.Jwt.JwtSecret, conf.Jwt.AccessTokerExpire, conf.Jwt.RefreshTokenExpire)

	noteRouter := NewRouter(notesRepo, logger, jwtServ)
	return &Server{
		router: noteRouter.router,
		port:   conf.Port,
	}
}

func (s *Server) Start() {
	s.router.Run(fmt.Sprintf(":%d", s.port))
}
