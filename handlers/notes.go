package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Создание новой заметки
func CreateNoteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, "CreateNoteHandler")
}

// Получение заметки по ID
func GetNoteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "GetNoteHandler")
}

// Изменение заметки по ID
func UpdateNoteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "UpdateNoteHandler")
}

// Удаление заметки по ID
func DeleteNoteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "DeleteNoteHandler")
}

// Получение списка всех заметок
func GetListNotesHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "GetListNotesHandler")
}
