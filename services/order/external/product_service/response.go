package productservice

import (
	"time"

	"github.com/google/uuid"
)

type GetProductByIDRespData struct {
	ID          uuid.UUID `json:"id"`
	ShopID      uuid.UUID `json:"shop_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetProductByIDResp struct {
	Data    GetProductByIDRespData `json:"data"`
	Success string                 `json:"success"`
}

type ErrorResponse struct {
	Metadata ErrorMetadata `json:"metadata"`
}

type ErrorMetadata struct {
	Path    string `json:"path"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
