package model

import (
	"time"
)

type Score struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username" gorm:"size:100;index"` // добавили ограничение + индекс
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
