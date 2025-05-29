package payload

type ReserveStocksReq struct {
	Stocks []ReserveStocksData `json:"stocks" binding:"required,dive"`
}

type ReserveStocksData struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type ReserveStocksResp struct {
	ProductID        string `json:"product_id"`
	WarehouseID      string `json:"warehouse_id"`
	ReservedQuantity int    `json:"reserved_quantity"`
}
