package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notes-manager/database"
	"notes-manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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

	go resetCache(fmt.Sprintf("notes/%d", note.AuthorId))
	ctx.JSON(http.StatusOK, gin.H{
		"note":    result,
		"message": "Заметка успешно добавлена",
	})

}

// Получение заметки по ID
func GetNoteHandler(ctx *gin.Context) {
	var note models.Note
	id := ctx.Param("id")
	filter := bson.M{"id": id}

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", 1))

	if err := collection.FindOne(ctx, filter).Decode(&note); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": note})
}

// Изменение заметки по ID
func UpdateNoteHandler(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorId := 1
	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", authorId))

	updateFilds := bson.M{}
	if note.Name != nil {
		updateFilds["name"] = note.Name
	}
	if note.Content != nil {
		updateFilds["content"] = note.Content
	}
	update := bson.M{"$set": updateFilds}

	id := ctx.Param("id")
	filter := bson.M{"id": id}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Заметка не найдена"})
		return
	}

	go resetCache(fmt.Sprintf("notes/%d", note.AuthorId))
	ctx.JSON(http.StatusOK, gin.H{"message": result})
}

// Удаление заметки по ID
func DeleteNoteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	filter := bson.M{"id": id}
	authorId := 1

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", authorId))

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error:": "запись не удалена"})
		return
	}

	go resetCache(fmt.Sprintf("notes/%d", authorId))
	ctx.JSON(http.StatusOK, gin.H{"message": result})
}

// Получение списка всех заметок
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

func resetCache(collection string) (int64, error) {
	return database.RedisClient.Del(collection).Result()
}
