package model

import "github.com/google/uuid"

type GetStockAvailablesByProduct struct {
	ProductID      uuid.UUID `json:"product_id"`
	AvailableStock int       `json:"available_stock"`
}
