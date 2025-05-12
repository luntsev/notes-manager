package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luntsev/notes-manager/auth/models"
	"github.com/luntsev/notes-manager/auth/pkg/jwt"
	"github.com/luntsev/notes-manager/auth/repository"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Неудалось создать токены"})
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
	var refreshToken models.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&refreshToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные обновления токена"})
		return
	}

	jwdData, err := h.jwtServ.Verify(refreshToken.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Ошибка при проверке токена обновления: %s", err.Error())})
		return
	}

	tokens, err := h.jwtServ.Create(jwdData)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Неудалось создать токены"})
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (h *authHandler) Update(ctx *gin.Context) {
	var updateData models.UpdateRequest

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные обновления учетной записи"})
		return
	}

	email := ctx.GetString("email")
	user, err := h.repo.GetByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Учетная запись пользователя не найдена"})
		h.logger.WriteError(fmt.Sprintf("unable find user account: %s", err.Error()))
		return
	}

	if updateData.Name != nil {
		user.Name = *updateData.Name
	}

	if updateData.Password != nil {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(*updateData.Password), 10)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось хешировать пароль"})
			h.logger.WriteError(fmt.Sprintf("Unable hashing password: %s", err.Error()))
			return
		}
		user.PassHash = string(hashedPass)
	}

	if updateData.BirthDay != nil {
		user.BirthDay = updateData.BirthDay.Time
	}

	err = h.repo.Update(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалосьобновить аккаунт пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь обновлен"})
}
