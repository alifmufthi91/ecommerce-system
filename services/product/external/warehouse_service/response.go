package warehouseservice

type GetStockAvailablesData struct {
	ProductID      string `json:"product_id"`
	AvailableStock int    `json:"available_stock"`
}

type GetStockAvailablesResp struct {
	Data    []GetStockAvailablesData `json:"data"`
	Success string                   `json:"success"`
}

type ErrorResponse struct {
	Metadata ErrorMetadata `json:"metadata"`
}

type ErrorMetadata struct {
	Path    string `json:"path"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
