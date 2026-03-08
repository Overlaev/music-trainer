package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Score struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

var db *gorm.DB

func main() {
	var err error
	// 2. Подключаемся к SQLite (файл создастся автоматически)
	db, err = gorm.Open(sqlite.Open("leaderboard.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	// Автоматическая миграция (создание таблицы)
	db.AutoMigrate(&Score{})

	// 3. Настраиваем Gin
	r := gin.Default()

	// Настраиваем CORS, чтобы фронтенд мог делать запросы к бэкенду
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Для MVP разрешаем запросы отовсюду
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// 4. Роуты (API Эндпоинты)
	api := r.Group("/api")
	{
		api.POST("/scores", saveScore)          // Сохранить результат
		api.GET("/leaderboard", getLeaderboard) // Получить топ
	}

	// 5. Запускаем сервер на порту 8080
	log.Println("Сервер запущен на http://localhost:8080")
	r.Run(":8080")
}

func saveScore(c *gin.Context) {
	var input Score

	// Парсим JSON из запроса фронтенда
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Если имя не передали, ставим "Аноним"
	if input.Username == "" {
		input.Username = "Аноним"
	}

	// Сохраняем в БД
	db.Create(&input)

	c.JSON(http.StatusOK, gin.H{"message": "Результат сохранен!", "data": input})
}

func getLeaderboard(c *gin.Context) {
	var topScores []Score

	// Ищем топ-10 результатов в базе, сортируем по очкам по убыванию
	if err := db.Order("score desc").Limit(10).Find(&topScores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить таблицу лидеров"})
		return
	}

	c.JSON(http.StatusOK, topScores)
}
