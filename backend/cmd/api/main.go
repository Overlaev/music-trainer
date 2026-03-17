package main

import (
	"log"

	"music-trainer/internal/model"

	"music-trainer/internal/server"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("leaderboard.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Миграция
	if err := db.AutoMigrate(&model.Score{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	router := server.NewRouter(db)

	log.Println("Сервер запущен → http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
