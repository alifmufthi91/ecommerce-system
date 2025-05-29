package handler

import (
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/httpresp"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/observ"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
)

// @Summary		Shop - Get Shops
// @Description	get all shops
// @Tags		Shop
// @Accept		json
// @Produce		json
// @Success		200	{object}	httpresp.Response{data=[]model.Shop}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/shops [get]
func (h *shopHandler) GetShops(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "shopHandler.GetShops")
	defer span.End()

	shops, err := h.shopService.GetShops(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, shops, nil)
}
