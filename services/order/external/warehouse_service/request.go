package warehouseservice

type ReserveStocksReq struct {
	Stocks []ReserveStocksReqData `json:"stocks"`
	Token  string                 `json:"-"`
}

type ReserveStocksReqData struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CommitReservesReq struct {
	Stocks []CommitReservesReqData `json:"stocks"`
	Token  string                  `json:"-"`
}

type CommitReservesReqData struct {
	ProductID   string `json:"product_id"`
	WarehouseID string `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
}

type RollbackReservesReq struct {
	Stocks []RollbackReservesReqData `json:"stocks"`
	Token  string                    `json:"-"`
}

type RollbackReservesReqData struct {
	ProductID   string `json:"product_id"`
	WarehouseID string `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
}
