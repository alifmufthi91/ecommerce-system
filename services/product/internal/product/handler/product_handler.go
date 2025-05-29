package handler

import (
	"strings"

	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/auth"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/httpresp"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/utils"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/payload"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
)

// @Summary		Product - Create Product
// @Description	create a new product
// @Tags		Product
// @Accept		json
// @Produce		json
// @param		request	body	payload.CreateProductReq	true	"create product request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/products [post]
func (h *productHandler) CreateProduct(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "productHandler.CreateProduct")
	defer span.End()

	var req payload.CreateProductReq
	if err := c.BindJSON(&req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	err := h.productService.CreateProduct(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, "success", nil)
}

// @Summary		Product - Get Products
// @Description	get all products
// @Tags		Product
// @Accept		json
// @Produce		json
// @Success		200	{object}	httpresp.Response{data=[]payload.GetProductsResp}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/products [get]
func (h *productHandler) GetProducts(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "productHandler.GetProducts")
	defer span.End()

	claims := auth.GetClaimsFromContext(c)
	if claims == nil {
		httpresp.HttpRespError(c, apperr.NewWithCode(apperr.CodeHTTPUnauthorized, `Not authorized`))
		return
	}

	products, err := h.productService.GetProducts(ctx, claims.Token)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, products, nil)
}

// @Summary		Product - Get Product By ID
// @Description	get product by ID
// @Tags		Product
// @Accept		json
// @Produce		json
// @Param		id	path	string	true	"product ID"
// @Success		200	{object}	httpresp.Response{data=model.Product}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Security	BearerAuth
// @Router		/products/{id} [get]
func (h *productHandler) GetProductByID(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "productHandler.GetProductByID")
	defer span.End()

	productID := c.Param("id")
	if productID == "" {
		httpresp.HttpRespError(c, apperr.NewWithCode(apperr.CodeHTTPBadRequest, "product ID is required"))
		return
	}

	product, err := h.productService.GetProductByID(ctx, productID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, product, nil)
}
