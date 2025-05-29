package model

import (
	"time"

	"github.com/google/uuid"
)

type StockLock struct {
	ID          uuid.UUID `json:"id" gorm:"column:id;primaryKey;default:uuid_generate_v4()"`
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
