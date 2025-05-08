package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"github.com/notes-manager/auth/configs"
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

	logger := logs.NewLogger(conf.LogLevel)

	router := NewAuthRouter()

	return &authServer{
		router: router.router,
		port:   conf.Port,
	}
}

func (s *authServer) Start() {
	s.router.Run(fmt.Sprintf(":%d", s.port))
}
