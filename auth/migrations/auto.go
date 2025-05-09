package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/luntsev/notes-manager/auth/configs"
	"github.com/luntsev/notes-manager/auth/models"
	"github.com/luntsev/notes-manager/auth/pkg/database"
	"github.com/luntsev/notes-manager/notes/pkg/logs"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("File \".env\" not foud: %s. Using default config.", err.Error())
	}
}

func main() {
	conf, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := logs.NewLogger(conf.LogLevel)

	db := database.NewPostgresDB(conf, logger)
	db.DB.AutoMigrate(&models.User{})

}
