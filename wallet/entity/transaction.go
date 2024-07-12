package entity

import (
	"time"
)

type Transaction struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID    int       `json:"sender_id"`
	RecipientID int       `json:"recipient_id"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
