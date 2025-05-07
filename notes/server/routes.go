package server

import (
	"notes-manager/handlers"
	"notes-manager/pkg/logs"
	"notes-manager/repository"

	"github.com/gin-gonic/gin"
)

type NotesHandler struct {
	*repository.NotesRepository
	*logs.Logger
}

func InitRoutes(repo *repository.NotesRepository, logger *logs.Logger) *gin.Engine {
	handlers := handlers.NewNotesHandler(repo, logger)

	router := gin.Default()

	router.POST("/note", handlers.Create)       // Создание заметки
	router.GET("/note/:id", handlers.Get)       // Получение заметки
	router.PUT("/note/:id", handlers.Update)    // Изменение заметки
	router.DELETE("/note/:id", handlers.Delete) // Удаление заметки
	router.GET("/notes", handlers.GetList)      // Получение списка заметок

	return router
}
