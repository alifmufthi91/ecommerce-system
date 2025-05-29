package payload

import "github.com/google/uuid"

type TransferStockReq struct {
	FromWarehouseID uuid.UUID `json:"from_warehouse_id" binding:"required"`
	ToWarehouseID   uuid.UUID `json:"to_warehouse_id" binding:"required"`
	ProductID       uuid.UUID `json:"product_id" binding:"required"`
	Quantity        int       `json:"quantity" binding:"required,min=1"`
}
