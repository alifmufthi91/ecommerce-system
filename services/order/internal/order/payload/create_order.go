package payload

import "github.com/google/uuid"

type CreateOrderReq struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	UserID    string    `json:"-"`
	Token     string    `json:"-"`
}
