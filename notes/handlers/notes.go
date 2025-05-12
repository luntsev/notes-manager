package handlers

import (
	"fmt"
	"net/http"

	"github.com/luntsev/notes-manager/auth/pkg/jwt"
	"github.com/luntsev/notes-manager/notes/models"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"github.com/luntsev/notes-manager/notes/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type notesHandler struct {
	repo    *repository.NotesRepository
	logger  *logs.Logger
	jwtServ *jwt.JWT
}

func NewNotesHandler(notesRepo *repository.NotesRepository, logs *logs.Logger, jwt *jwt.JWT) *notesHandler {
	return &notesHandler{
		repo:    notesRepo,
		logger:  logs,
		jwtServ: jwt,
	}
}

// Создание новой заметки
func (h *notesHandler) Create(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.WriteWarn(fmt.Sprintf("Bad request to create note: %s", err.Error()))
		return
	}

	note.Id = uuid.NewString()
	note.AuthorId = ctx.GetUint("id")

	result, err := h.repo.Create(&note, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"note":    result,
		"message": "Заметка успешно добавлена",
	})
}

// Получение заметки по ID
func (h *notesHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	authorId := ctx.GetUint("id")

	note, err := h.repo.Read(id, authorId, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": note})
}

// Изменение заметки по ID
func (h *notesHandler) Update(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.logger.WriteWarn(fmt.Sprintf("Bad request to update note: %s", err.Error()))
		return
	}

	note.AuthorId = ctx.GetUint("id")
	note.Id = ctx.Param("id")
	err := h.repo.Update(&note, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Заметка обновлена"})
}

// Удаление заметки по ID
func (h *notesHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	authorId := ctx.GetUint("id")

	delCount, err := h.repo.Delete(id, authorId, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}

	if delCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": "запись не удалена"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Удалено %d записей", delCount)})
}

// Получение списка всех заметок
func (h *notesHandler) GetList(ctx *gin.Context) {
	authorId := ctx.GetUint("id")

	notes, err := h.repo.List(authorId, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": notes})

}
