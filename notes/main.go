package main

import (
	"log"
	"notes-manager/configs"
	"notes-manager/pkg/database"
	"notes-manager/repository"
	"notes-manager/server"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("File \".env\" not foud: %s. Using default config.", err.Error())
	}
}

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	notesDb := database.NewMongoDB(config)
	notesRepo := repository.NewNotesRepository(config, notesDb)

	if err != nil {
		panic(err.Error())
	}
	server.StartServer()
}
