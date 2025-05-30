package handler

import (
	"strings"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/httpresp"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/utils"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/payload"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
)

// @Summary		Stock - Get Stocks
// @Description	get all stocks
// @Tags		Stock
// @Accept		json
// @Produce		json
// @Param		request	query	payload.GetStocksReq	false	"get stocks request query parameters"
// @Success		200	{object}	httpresp.Response{data=[]model.WarehouseStock}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/stocks [get]
func (h *stockHandler) GetStocks(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "stockHandler.GetStocks")
	defer span.End()

	var req payload.GetStocksReq
	if err := c.BindQuery(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	stocks, err := h.stockService.GetStocks(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, stocks, nil)
}

// @Summary		Stock - Transfer Stock
// @Description	transfer stock between warehouses
// @Tags		Stock
// @Accept		json
// @Produce		json
// @Param		request	body	payload.TransferStockReq	true	"transfer stock request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/stocks/transfer [post]
func (h *stockHandler) TransferStock(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "stockHandler.TransferStock")
	defer span.End()

	var req payload.TransferStockReq
	if err := c.BindJSON(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	if err := h.stockService.TransferStock(ctx, req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, "success", nil)
}

// @Summary		Stock - Get Available Stocks By Product
// @Description	get available stocks by product
// @Tags		Stock
// @Accept		json
// @Produce		json
// @Param		request	query	payload.GetStockAvailablesByProductReq	false	"get available stocks by product request query parameters"
// @Success		200	{object}	httpresp.Response{data=[]model.GetStockAvailablesByProduct}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/stocks/availables [get]
func (h *stockHandler) GetAvailableStocksByProduct(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "stockHandler.GetAvailableStocksByProduct")
	defer span.End()

	var req payload.GetStockAvailablesByProductReq
	if err := c.BindQuery(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	stocks, err := h.stockService.GetStockAvailablesByProduct(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, stocks, nil)
}

// @Summary		Stock - Reserve Stocks
// @Description	reserve stocks for an order
// @Tags		Stock
// @Accept		json
// @Produce		json
// @Param		request	body	payload.ReserveStocksReq	true	"reserve stocks request body"
// @Success		200	{object}	httpresp.Response{data=[]payload.ReserveStocksResp}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/stocks/reserve [post]
func (h *stockHandler) ReserveStocks(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "stockHandler.ReserveStocks")
	defer span.End()

	var req payload.ReserveStocksReq
	if err := c.BindJSON(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	reservedStocks, err := h.stockService.ReserveStocks(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, reservedStocks, nil)
}

// @Summary		Stock - Rollback Reserves
// @Description	rollback reserved stocks
// @Tags		Stock
// @Accept		json
// @Produce		json
// @Param		request	body	payload.RollbackReservesReq	true	"rollback reserves request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/stocks/rollback-reserves [post]
func (h *stockHandler) RollbackReserves(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "stockHandler.RollbackReserveStocks")
	defer span.End()

	var req payload.RollbackReservesReq
	if err := c.BindJSON(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	if err := h.stockService.RollbackReserves(ctx, req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, "success", nil)
}

// @Summary		Stock - Commit Reserves
// @Description	commit reserved stocks
// @Tags		Stock
// @Accept		json
// @Produce		json
// @Param		request	body	payload.CommitReservesReq	true	"commit reserves request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/stocks/commit [post]
func (h *stockHandler) CommitReserves(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "stockHandler.CommitReserves")
	defer span.End()

	var req payload.CommitReservesReq
	if err := c.BindJSON(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	if err := h.stockService.CommitReserves(ctx, req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, "success", nil)
}

// @Summary		Stock - Create Stock
// @Description	create a new stock
// @Tags		Stock
// @Accept		json
// @Produce		json
// @Param		request	body	payload.CreateStockReq	true	"create stock request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/stocks [post]
func (h *stockHandler) CreateStock(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "stockHandler.CreateStock")
	defer span.End()

	var req payload.CreateStockReq
	if err := c.BindJSON(&req); err != nil {
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}
	if err := h.stockService.CreateStock(ctx, req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}
	httpresp.HttpRespSuccess(c, "success", nil)
}
