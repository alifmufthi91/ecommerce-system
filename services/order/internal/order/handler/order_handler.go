package handler

import (
	"strings"

	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/payload"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/auth"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/httpresp"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
)

// @Summary		Order - Get Orders
// @Description	get all orders
// @Tags		Order
// @Accept		json
// @Produce		json
// @Param		request	query	payload.GetOrdersReq	false	"get orders request query parameters"
// @Success		200	{object}	httpresp.Response{data=[]model.Order}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/orders [get]
func (h *orderHandler) GetOrders(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "orderHandler.GetOrders")
	defer span.End()

	var req payload.GetOrdersReq
	if err := c.BindQuery(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	orders, err := h.orderService.GetOrders(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, orders, nil)
}

// @Summary		Order - Create Order
// @Description	create a new order
// @Tags		Order
// @Accept		json
// @Produce		json
// @Param		request	body	payload.CreateOrderReq	true	"create order request body"
// @Success		200	{object}	httpresp.Response{data=model.Order}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/orders [post]
func (h *orderHandler) CreateOrder(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "orderHandler.CreateOrder")
	defer span.End()

	var req payload.CreateOrderReq
	if err := c.BindJSON(&req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	claims := auth.GetClaimsFromContext(c)

	req.UserID = claims.UserID
	req.Token = claims.Token
	order, err := h.orderService.CreateOrder(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, order, nil)
}

// @Summary		Order - Complete Order
// @Description	complete an order by ID
// @Tags		Order
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"Order ID"
// @Success		200	{object}	httpresp.Response{data=model.Order}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/orders/{id}/complete [patch]
func (h *orderHandler) CompleteOrder(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "orderHandler.CompleteOrder")
	defer span.End()

	id := c.Param("id")
	if id == "" {
		span.SetStatus(codes.Error, "order ID is required")
		httpresp.HttpRespError(c, apperr.NewWithCode(apperr.CodeHTTPBadRequest, "order ID is required"))
		return
	}

	claims := auth.GetClaimsFromContext(c)

	var req payload.CompleteOrderReq
	req.OrderID = id
	req.Token = claims.Token
	order, err := h.orderService.CompleteOrder(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, order, nil)
}
