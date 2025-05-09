package server

import (
	"github.com/luntsev/notes-manager/auth/handlers"
	"github.com/luntsev/notes-manager/auth/pkg/jwt"
	"github.com/luntsev/notes-manager/auth/repository"
	"github.com/luntsev/notes-manager/notes/pkg/logs"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	router *gin.Engine
}

func NewAuthRouter(repo *repository.AuthRepository, logger *logs.Logger, jwtServ *jwt.JWT) *authRouter {
	router := gin.Default()

	handler := handlers.NewAuthHandler(repo, logger, jwtServ)

	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)
	router.PUT("/auth/update", handler.Update)
	router.GET("/auth", handler.GetUser)
	router.POST("/auth/refresh", handler.RefreshToken)

	return &authRouter{router: router}
}
