package payload

type CommitReservesReq struct {
	Stocks []CommitReservesData `json:"stocks" binding:"required,dive"`
}
type CommitReservesData struct {
	ProductID   string `json:"product_id" binding:"required"`
	WarehouseID string `json:"warehouse_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,gt=0"`
}
