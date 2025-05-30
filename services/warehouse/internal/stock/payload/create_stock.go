package payload

import "github.com/google/uuid"

type CreateStockReq struct {
	ProductID   uuid.UUID `json:"product_id" validate:"required,uuid"`
	WarehouseID uuid.UUID `json:"warehouse_id" validate:"required,uuid"`
	Quantity    int       `json:"quantity" validate:"required,min=1"`
}
