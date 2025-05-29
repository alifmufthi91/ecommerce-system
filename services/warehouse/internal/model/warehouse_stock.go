package model

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseStock struct {
	ID          uuid.UUID `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Reserved    int       `json:"reserved"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
