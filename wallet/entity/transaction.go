package entity

import (
	"time"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID           int       `json:"id"`
	FromWalletID int       `json:"from_wallet_id"`
	ToWalletID   int       `json:"to_wallet_id"`
	Amount       float64   `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
}
