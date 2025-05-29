package warehouseservice

import "github.com/google/uuid"

type ReserveStocksRespData struct {
	WarehouseID      uuid.UUID `json:"warehouse_id"`
	ProductID        uuid.UUID `json:"product_id"`
	ReservedQuantity int       `json:"reserved_quantity"`
}

type ReserveStocksResp struct {
	Data    []ReserveStocksRespData `json:"data"`
	Success string                  `json:"success"`
}

type CommitReservesResp struct {
	Data    any    `json:"data"`
	Success string `json:"success"`
}

type RollbackReservesResp struct {
	Data    any    `json:"data"`
	Success string `json:"success"`
}

type ErrorResponse struct {
	Metadata ErrorMetadata `json:"metadata"`
}

type ErrorMetadata struct {
	Path    string `json:"path"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
