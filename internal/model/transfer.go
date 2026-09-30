package model

import (
	"fmt"
	"time"
)

type Transfer struct {
	TransactionID string    `json:"transaction_id" binding:"required" bson:"transaction_id"`
	Sender        string    `json:"sender" binding:"required" bson:"sender"`
	Receiver      string    `json:"receiver" binding:"required" bson:"receiver"`
	Amount        int64     `json:"amount" binding:"required" bson:"amount"`
	Currency      string    `json:"currency" binding:"required" bson:"currency"`
	OccurredAt    time.Time `json:"occurred_at" binding:"required" bson:"occurred_at"`
	ReceivedAt    time.Time `json:"-" bson:"received_at"`
}

func (t *Transfer) Validate() error {
	if t.TransactionID == "" {
		return fmt.Errorf("transaction_id is required")
	}
	if t.Sender == "" {
		return fmt.Errorf("sender is required")
	}
	if t.Receiver == "" {
		return fmt.Errorf("receiver is required")
	}
	if t.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if t.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if t.OccurredAt.IsZero() {
		return fmt.Errorf("occurred_at is required")
	}
	return nil
}