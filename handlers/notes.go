package handlers

import (
	"fmt"
	"net/http"
	"notes-manager/database"
	"notes-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Создание новой заметки
func CreateNoteHandler(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note.Id = uuid.NewString()
	note.AuthorId = 1

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", note.AuthorId))
	result, err := collection.InsertOne(ctx, note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"note":    result,
		"message": "Заметка успешно добавлена",
	})

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
	id := ctx.Param("id")
	filter := bson.M{"id": id}

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", 1))

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": "запись не удалена"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": result})
}

// Получение списка всех заметок
func GetListNotesHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "GetListNotesHandler")
}
