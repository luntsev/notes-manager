package server

import (
	"github.com/notes-manager/auth/handlers"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	router *gin.Engine
}

func NewAuthRouter() *authRouter {
	router := gin.Default()

	handler := handlers.NewAuthHandler()

	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)
	router.PUT("/auth/update", handler.Update)
	router.GET("/auth", handler.GetUser)
	router.POST("/auth/refresh", handler.RefreshToken)

	return &authRouter{router: router}
}
