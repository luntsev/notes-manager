package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/luntsev/notes-manager/auth/configs"
	"github.com/luntsev/notes-manager/auth/pkg/database"
	"github.com/luntsev/notes-manager/auth/pkg/jwt"
	"github.com/luntsev/notes-manager/auth/pkg/middleware"
	"github.com/luntsev/notes-manager/auth/repository"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
)

type authServer struct {
	router *gin.Engine
	port   int
}

func NewServer() *authServer {
	conf, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	jwtServ := jwt.NewJWT(conf.Jwt.JwtSecret, conf.Jwt.AccessTokerExpire, conf.Jwt.RefreshTokenExpire)
	logger := logs.NewLogger(conf.LogLevel)

	db := database.NewPostgresDB(conf, logger)
	repo := repository.NewAuthRepository(db, logger)

	router := NewAuthRouter(repo, logger, jwtServ)
	router.router.Use(middleware.IsAuth(&conf.Jwt))

	return &authServer{
		router: router.router,
		port:   conf.Port,
	}
}

func (s *authServer) Start() {
	s.router.Run(fmt.Sprintf(":%d", s.port))
}
