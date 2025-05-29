package payload

type CompleteOrderReq struct {
	OrderID string `json:"-"`
	Token   string `json:"-"`
}
