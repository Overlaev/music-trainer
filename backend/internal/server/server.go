package server

import (
	"music-trainer/internal/handler"
	"music-trainer/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS — для продакшена лучше указать конкретные домены
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ← потом поменять!
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 60 * 60, // 12 часов
	}))

	// Внедряем зависимости
	scoreRepo := repository.NewScoreRepository(db) // ← обратите внимание на название пакета
	scoreHandler := handler.NewScoreHandler(scoreRepo)

	// Группа API
	api := r.Group("/api")
	{
		api.POST("/scores", scoreHandler.SaveScore)
		api.GET("/leaderboard", scoreHandler.GetLeaderboard)
	}

	return r
}
