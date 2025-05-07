package main

import (
	"github.com/joho/godotenv"
	"log"
	"notes-manager/server"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("File \".env\" not foud: %s. Using default config.", err.Error())
	}
}

func main() {
	router := server.InitServer()
	server.StartServer(router)
}
