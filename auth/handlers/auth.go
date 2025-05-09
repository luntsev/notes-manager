package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luntsev/notes-manager/auth/models"
	"github.com/luntsev/notes-manager/auth/repository"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	repo   *repository.AuthRepositury
	logger *logs.Logger
}

func NewAuthHandler(repository *repository.AuthRepositury, logger *logs.Logger) *authHandler {
	return &authHandler{
		repo:   repository,
		logger: logger,
	}
}

func (h *authHandler) Register(ctx *gin.Context) {
	var registerData models.RegisterData

	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные регистрации"})
		h.logger.WriteWarn("Wrong register data in request")
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
		h.logger.WriteError(fmt.Sprintf("Unable hashing password: %s", err.Error()))
		return
	}

	user := &models.User{
		Email:    registerData.Email,
		PassHash: string(hashedPass),
	}

	if err := h.repo.Create(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалосб создать пользователя"})
		h.logger.WriteError(fmt.Sprintf("Unable write user record: %s", err.Error()))
	}

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
