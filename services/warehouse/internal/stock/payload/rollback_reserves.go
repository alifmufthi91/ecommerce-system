package payload

type RollbackReservesReq struct {
	Stocks []RollbackReservesData `json:"stocks" binding:"required,dive"`
}
type RollbackReservesData struct {
	ProductID   string `json:"product_id" binding:"required"`
	WarehouseID string `json:"warehouse_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,gt=0"`
}
