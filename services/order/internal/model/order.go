package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID `json:"user_id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status" gorm:"default:'pending'"` // e.g., pending, completed, cancelled
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
