package server

import (
	"github.com/luntsev/notes-manager/notes/handlers"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"github.com/luntsev/notes-manager/notes/repository"

	"github.com/gin-gonic/gin"
)

type notesRouter struct {
	router *gin.Engine
}

func NewRouter(repo *repository.NotesRepository, logger *logs.Logger) *notesRouter {
	handlers := handlers.NewNotesHandler(repo, logger)

	router := gin.Default()

	router.POST("/note", handlers.Create)       // Создание заметки
	router.GET("/note/:id", handlers.Get)       // Получение заметки
	router.PUT("/note/:id", handlers.Update)    // Изменение заметки
	router.DELETE("/note/:id", handlers.Delete) // Удаление заметки
	router.GET("/notes", handlers.GetList)      // Получение списка заметок

	return &notesRouter{router: router}
}
