package productservice

type GetProductByIDReq struct {
	ProductID string `json:"product_id"`
	Token     string `json:"-"`
}
