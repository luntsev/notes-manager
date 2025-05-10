package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luntsev/notes-manager/auth/models"
	"github.com/luntsev/notes-manager/auth/pkg/jwt"
	"github.com/luntsev/notes-manager/auth/repository"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	repo    *repository.AuthRepository
	logger  *logs.Logger
	jwtServ *jwt.JWT
}

func NewAuthHandler(repository *repository.AuthRepository, logger *logs.Logger, jwt *jwt.JWT) *authHandler {
	return &authHandler{
		repo:    repository,
		logger:  logger,
		jwtServ: jwt,
	}
}

func (h *authHandler) Register(ctx *gin.Context) {
	var regData models.RegisterRequest

	if err := ctx.ShouldBindJSON(&regData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные регистрации"})
		h.logger.WriteWarn("Wrong register data in request")
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(regData.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
		h.logger.WriteError(fmt.Sprintf("Unable hashing password: %s", err.Error()))
		return
	}

	user := &models.User{
		Email:    regData.Email,
		Name:     regData.Name,
		PassHash: string(hashedPass),
		BirthDay: regData.BirthDay.Time,
	}

	if err := h.repo.Create(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалосб создать пользователя"})
		h.logger.WriteError(fmt.Sprintf("Unable write user record: %s", err.Error()))
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Пользователь зарегистрирован"})
}

func (h *authHandler) Login(ctx *gin.Context) {
	var regData models.LoginRequest

	if err := ctx.ShouldBindJSON(&regData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные авторизации"})
		h.logger.WriteWarn("Wrong login data in request")
		return
	}

	user, err := h.repo.GetByEmail(regData.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный email или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(regData.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный email или пароль"})
		h.logger.WriteInfo(fmt.Sprintf("Wrong password for account: %s, fromt IP: %s", regData.Email, ctx.ClientIP()))
		return
	}

	tokens, err := h.jwtServ.Create(&jwt.JWTData{Email: regData.Email})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable create tokens"})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (h *authHandler) GetUser(ctx *gin.Context) {
	email := ctx.GetString("email")
	if email == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Авторизация не выполнена"})
		return
	}

	user, err := h.repo.GetByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь не найден"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": user})
}

func (h *authHandler) RefreshToken(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Токен обновлен"})
}

func (h *authHandler) Update(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь обновлен"})
}
