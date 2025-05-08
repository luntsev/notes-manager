package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type authServer struct {
	router *gin.Engine
}

func NewServer() *authServer {
	router := NewAuthRouter()

	return &authServer{
		router: router.router,
	}
}

func (s *authServer) Start(port int) {
	if port <= 0 {
		log.Fatalf("Invalid port: %d", port)
	}
	s.router.Run(fmt.Sprintf(":%d", port))
}
