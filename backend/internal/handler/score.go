package handler

import (
	"net/http"

	"music-trainer/internal/model"
	"music-trainer/internal/repository"

	"github.com/gin-gonic/gin"
)

type ScoreHandler struct {
	repo repository.ScoreRepository
}

func NewScoreHandler(repo repository.ScoreRepository) *ScoreHandler {
	return &ScoreHandler{repo: repo}
}

func (h *ScoreHandler) SaveScore(c *gin.Context) {
	var input model.Score
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Username == "" {
		input.Username = "Аноним"
	}

	if err := h.repo.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить результат"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Результат сохранён",
		"data":    input,
	})
}

func (h *ScoreHandler) GetLeaderboard(c *gin.Context) {
	scores, err := h.repo.GetTop(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить таблицу лидеров"})
		return
	}

	c.JSON(http.StatusOK, scores)
}
