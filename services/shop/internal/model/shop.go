package model

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID        uuid.UUID `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
