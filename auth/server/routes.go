package server

import (
	"github.com/luntsev/notes-manager/auth/handlers"
	"github.com/luntsev/notes-manager/auth/repository"
	"github.com/luntsev/notes-manager/notes/pkg/logs"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	router *gin.Engine
}

func NewAuthRouter(repo *repository.AuthRepositury, logger *logs.Logger) *authRouter {
	router := gin.Default()

	handler := handlers.NewAuthHandler(repo, logger)

	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)
	router.PUT("/auth/update", handler.Update)
	router.GET("/auth", handler.GetUser)
	router.POST("/auth/refresh", handler.RefreshToken)

	return &authRouter{router: router}
}
