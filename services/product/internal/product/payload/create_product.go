package payload

import "github.com/google/uuid"

type CreateProductReq struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	ShopID      uuid.UUID `json:"shop_id" binding:"required"`
}
