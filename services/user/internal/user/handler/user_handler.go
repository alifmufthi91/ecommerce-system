package handler

import (
	"strings"

	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/httpresp"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/utils"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/payload"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
)

// @Summary		User - Register User
// @Description	create a new user account
// @Tags		User
// @Accept		json
// @Produce		json
// @param		request	body	payload.RegisterUserReq	true	"register user request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Router		/users [post]
func (h *userHandler) Register(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "userHandler.Register")
	defer span.End()

	var req payload.RegisterUserReq
	if err := c.BindJSON(&req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	err := h.userService.RegisterUser(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, "success", nil)
}

// @Summary		User - Login User
// @Description	user login to get access token
// @Tags		User
// @Accept		json
// @Produce		json
// @param		request	body	payload.LoginUserReq	true	"login user request body"
// @Success		200	{object}	httpresp.Response{data=string}
// @Failure		400	{object}	httpresp.HTTPErrResp
// @Failure		404	{object}	httpresp.HTTPErrResp
// @Failure		500	{object}	httpresp.HTTPErrResp
// @Router		/users/login [post]
func (h *userHandler) Login(c *gin.Context) {
	ctx, span := observ.GetTracer().Start(c.Request.Context(), "userHandler.Login")
	defer span.End()

	var req payload.LoginUserReq
	if err := c.BindJSON(&req); err != nil {
		span.SetStatus(codes.Error, err.Error())
		errResp := strings.Join(utils.ParseBindErrors(err), "; ")
		httpresp.HttpRespError(c, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, errResp))
		return
	}

	token, err := h.userService.LoginUser(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, token, nil)
}
