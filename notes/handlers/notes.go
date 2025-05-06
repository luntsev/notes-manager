package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notes-manager/configs"
	"notes-manager/models"
	"notes-manager/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type NotesHandler struct {
	repo     *repository.NotesRepository
	logLevel string
}

func NewNotesHandler(notesRepo *repository.NotesRepository, conf configs.Config) *NotesHandler {
	return &NotesHandler{
		repo:     notesRepo,
		logLevel: conf.LogLevel,
	}
}

// Создание новой заметки
func (h *NotesHandler) Create(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.writeToLog(fmt.Sprint("Bad request to create note:", err.Error()))
		return
	}

	note.Id = uuid.NewString()
	note.AuthorId = 1

	result, err := h.repo.Create(&note, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.writeToLog(fmt.Sprint("Unable to create note", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"note":    result,
		"message": "Заметка успешно добавлена",
	})
	if h.logLevel == configs.DebugLog {
		log.Println("Note is Created")
	}
}

// Получение заметки по ID
func (h *NotesHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	var authorId uint = 1

	note, err := h.repo.Read(id, authorId, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		h.writeToLog("Note was not found when reading")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": note})
}

// Изменение заметки по ID
func (h *NotesHandler) Update(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.writeToLog(fmt.Sprint("Bad request to update note:", err.Error()))
		return
	}

	note.AuthorId = 1
	note.Id = ctx.Param("id")
	err := h.repo.Update(&note, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.writeToLog(fmt.Sprint("Note was not update:", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Заметка обновлена"})
	h.writeToLog("Note was update")
}

// Удаление заметки по ID
func (h *NotesHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var authorId uint = 1

	delCount, err := h.repo.Delete(id, authorId, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		h.writeToLog(fmt.Sprint("Error when delete note:", err.Error()))
		return
	}

	if delCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": "запись не удалена"})
		h.writeToLog(fmt.Sprint("Note was not delete:", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Удалено %d записей", delCount)})
	h.writeToLog("Note was delete")
}

// Получение списка всех заметок
func (h *NotesHandler) GetList(ctx *gin.Context) {
	var authorId uint = 1

	notes, err := h.repo.List(authorId, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetListNotesHandler(ctx *gin.Context) {
	authorId := 1
	var notes []models.Note

	cache, err := database.RedisClient.Get(fmt.Sprintf("notes/%d", authorId)).Result()
	if err == redis.Nil {
		log.Println("Кэш не найдей загружаем из БД")
		collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", authorId))

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var note models.Note

		for cursor.Next(ctx) {
			err := cursor.Decode(&note)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			notes = append(notes, note)
		}

		if len(notes) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"message": "Заметок не найдено"})
			return
		} else {
			go func() {
				err := recordToCache(notes, authorId)
				if err != nil {
					log.Println("При кешировании данных возникла ошибка:", err.Error())
				}
			}()
		}

	} else {
		log.Println("Кэш найдей загружаем из кэша")
		json.Unmarshal([]byte(cache), &notes)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": notes})
}

func recordToCache(notes []models.Note, authorId int) error {
	notesJson, err := json.Marshal(notes)
	if err != nil {
		log.Println("При кешировании данных возникла ошибка:", err.Error())
		return err
	}
	if err = database.RedisClient.Set(fmt.Sprintf("notes/%d", authorId), string(notesJson), time.Minute*1440).Err(); err != nil {
		return err
	}
	return nil
}

func (h *NotesHandler) writeToLog(msg string) {
	if h.logLevel == configs.DebugLog {
		log.Println(msg)
	}
}
