package payload

import (
	"time"

	"github.com/google/uuid"
)

type GetProductsResp struct {
	ID             uuid.UUID `json:"id"`
	ShopID         uuid.UUID `json:"shop_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	AvailableStock int       `json:"available_stock"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
