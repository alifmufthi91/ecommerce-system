package payload

type GetStocksReq struct {
	WarehouseIDIN []string `form:"warehouse_id_in" binding:"omitempty"`
	ProductIDIN   []string `form:"product_id_in" binding:"omitempty"`
}
