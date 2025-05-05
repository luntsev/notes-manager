package server

import (
	"log"
	"notes-manager/database"
	"notes-manager/envs"
)

func InitServer() {
	if err := envs.LoadEnvs(); err != nil {
		log.Fatal("Не удалось загрузить переменные окружения:", err.Error())
	} else {
		log.Println("Переменные окружения успешно загруженны")
	}

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
