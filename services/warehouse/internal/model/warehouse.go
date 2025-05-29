package model

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID        uuid.UUID `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Status    string    `json:"status"` // e.g., "active", "inactive"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
