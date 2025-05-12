package server

import (
	"github.com/luntsev/notes-manager/auth/pkg/jwt"
	"github.com/luntsev/notes-manager/auth/pkg/middleware"
	"github.com/luntsev/notes-manager/notes/handlers"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
	"github.com/luntsev/notes-manager/notes/repository"

	"github.com/gin-gonic/gin"
)

type notesRouter struct {
	router *gin.Engine
}

func NewRouter(repo *repository.NotesRepository, logger *logs.Logger, jwtServ *jwt.JWT) *notesRouter {
	handlers := handlers.NewNotesHandler(repo, logger, jwtServ)

	router := gin.Default()

	router.POST("/", middleware.IsAuth(jwtServ), handlers.Create)      // Создание заметки
	router.GET("/:id", middleware.IsAuth(jwtServ), handlers.Get)       // Получение заметки
	router.PUT("/:id", middleware.IsAuth(jwtServ), handlers.Update)    // Изменение заметки
	router.DELETE("/:id", middleware.IsAuth(jwtServ), handlers.Delete) // Удаление заметки
	router.GET("/list", middleware.IsAuth(jwtServ), handlers.GetList)  // Получение списка заметок

	return &notesRouter{router: router}
}
