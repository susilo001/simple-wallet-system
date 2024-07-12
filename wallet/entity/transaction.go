package entity

import (
	"time"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	RecipientID int       `json:"recipient_id"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
