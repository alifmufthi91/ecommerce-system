package payload

import "time"

type GetOrdersReq struct {
	UserIDIN      []string  `form:"user_id_in" binding:"omitempty"`
	ProductIDIN   []string  `form:"product_id_in" binding:"omitempty"`
	StatusIN      []string  `form:"status_in" binding:"omitempty"`
	ExpiresBefore time.Time `form:"expires_before" binding:"omitempty"`
}
