package warehouseservice

import "net/url"

type GetStockAvailablesReq struct {
	ProductIDIN []string
	Token       string
}

func (r *GetStockAvailablesReq) ToQueryParams() url.Values {
	queryParams := make(url.Values)

	if len(r.ProductIDIN) > 0 {
		for _, id := range r.ProductIDIN {
			queryParams.Add("product_id_in", id)
		}
	}

	return queryParams
}
