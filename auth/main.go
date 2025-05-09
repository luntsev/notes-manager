package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/luntsev/notes-manager/auth/server"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("File \".env\" not foud: %s. Using default config.", err.Error())
	}
}

func main() {
	noteServer := server.NewServer()
	noteServer.Start()
}
