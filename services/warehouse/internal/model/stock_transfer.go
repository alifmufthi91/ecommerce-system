package model

import (
	"time"

	"github.com/google/uuid"
)

type StockTransfer struct {
	ID              uuid.UUID `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	FromWarehouseID uuid.UUID `json:"from_warehouse_id"`
	ToWarehouseID   uuid.UUID `json:"to_warehouse_id"`
	ProductID       uuid.UUID `json:"product_id"`
	Quantity        int       `json:"quantity"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
