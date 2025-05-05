package server

import (
	"log"
	"notes-manager/database"
	"notes-manager/envs"
)

func InitServer() {
	envs.LoadEnvs()

	if err := database.InitDatabase(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err.Error())
	} else {
		log.Println("Выполнено подключение к базе данных")
	}

	if err := database.InitRedis(); err != nil {
		log.Fatal("Не удалось подключиться к Redis:", err.Error())
	} else {
		log.Println("Выполнено подключение к Redis")
	}
}

func StartServer() {
	InitRoutes()
}
