package handlers

import (
	"fmt"
	"net/http"
	"notes-manager/models"
	"notes-manager/pkg/logs"
	"notes-manager/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type notesHandler struct {
	repo   *repository.NotesRepository
	logger *logs.Logger
}

func NewNotesHandler(notesRepo *repository.NotesRepository, logs *logs.Logger) *notesHandler {
	return &notesHandler{
		repo:   notesRepo,
		logger: logs,
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
	note.AuthorId = 1

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
	var authorId uint = 1

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

	note.AuthorId = 1
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
	var authorId uint = 1

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
	var authorId uint = 1

	notes, err := h.repo.List(authorId, ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": notes})

}
