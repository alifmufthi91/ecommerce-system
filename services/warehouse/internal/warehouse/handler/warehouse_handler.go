package handler

import (
	"fmt"
	"strings"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/httpresp"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/utils"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/payload"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

// @Summary		Warehouse - Get Warehouses
// @Description	get all warehouses
// @Tags		Warehouse
// @Accept		json
// @Produce		json
// @Success		200	{object}	httpresp.Response{data=[]model.Warehouse}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/warehouses [get]
func (h *warehouseHandler) GetWarehouses(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "warehouseHandler.GetWarehouses")
	defer span.End()

	warehouses, err := h.warehouseService.GetWarehouses(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, warehouses, nil)
}

// @Summary		Warehouse - Update Warehouse
// @Description	update warehouse status
// @Tags		Warehouse
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"warehouse ID"
// @Param		request	body	payload.UpdateWarehouseReq	true	"update warehouse request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/warehouses/{id} [put]
func (h *warehouseHandler) UpdateWarehouse(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "warehouseHandler.UpdateWarehouse")
	defer span.End()

	id := c.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	var req payload.UpdateWarehouseReq
	if err := c.ShouldBind(&req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		fmt.Println(errResp)
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	req.ID = parsedId
	if err := h.warehouseService.UpdateWarehouse(ctx, req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, "success", nil)
}
