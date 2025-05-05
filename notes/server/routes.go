package server

import (
	"notes-manager/handlers"

	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	router := gin.Default()

	router.POST("/note", handlers.CreateNoteHandler)       // Создание заметки
	router.GET("/note/:id", handlers.GetNoteHandler)       // Получение заметки
	router.PUT("/note/:id", handlers.UpdateNoteHandler)    // Изменение заметки
	router.DELETE("/note/:id", handlers.DeleteNoteHandler) // Удаление заметки
	router.GET("/notes", handlers.GetListNotesHandler)     // Получение списка заметок

	router.Run(":9100")
}
