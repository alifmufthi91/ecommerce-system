package payload

type GetStockAvailablesByProductReq struct {
	ProductIDIN []string `form:"product_id_in" binding:"omitempty"`
}
