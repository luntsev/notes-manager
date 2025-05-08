package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct{}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}

func (h *authHandler) Register(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "Пользователь зарегистрирован"})
}

func (h *authHandler) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь авторизован"})
}

func (h *authHandler) GetUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Информация о пользователе получена"})
}

func (h *authHandler) RefreshToken(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Токен обновлен"})
}

func (h *authHandler) Update(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь обновлен"})
}
