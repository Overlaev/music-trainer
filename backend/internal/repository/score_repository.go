package repository

import (
	"music-trainer/internal/model"

	"gorm.io/gorm"
)

type ScoreRepository interface {
	Create(score *model.Score) error
	GetTop(limit int) ([]model.Score, error)
}

type scoreRepository struct {
	db *gorm.DB
}

func NewScoreRepository(db *gorm.DB) ScoreRepository {
	return &scoreRepository{db: db}
}

func (r *scoreRepository) Create(score *model.Score) error {
	return r.db.Create(score).Error
}

func (r *scoreRepository) GetTop(limit int) ([]model.Score, error) {
	var scores []model.Score
	err := r.db.
		Order("score desc").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}
