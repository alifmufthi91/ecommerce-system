package middleware

import (
	"strings"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/constant"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/auth"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/httpresp"
	"github.com/gin-gonic/gin"
)

func JwtMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		secretKey := cfg.Token.JWTSecret
		staticToken := cfg.Token.JWTStatic

		// token claims
		claims := &auth.CustomClaims{}
		headerToken, err := ParseTokenFromHeader(cfg, ctx)
		if err != nil {
			httpresp.HttpRespError(ctx, err)
			return
		}

		if headerToken == staticToken {
			// if static token, set claims with static token
			claims.UserEmail = constant.JWTStaticAdminEmail
			ctx.Set(constant.XUserEmail, claims.UserEmail)
			ctx.Set(auth.ContextClaimKey, claims)
			ctx.Next()
			return
		}

		claims, err = auth.ParseToken(secretKey, headerToken)
		if err != nil {
			httpresp.HttpRespError(ctx, apperr.NewWithCode(apperr.CodeHTTPUnauthorized, err.Error()))
			return
		}
		ctx.Set(constant.XUserEmail, claims.UserEmail)
		ctx.Set(auth.ContextClaimKey, claims)
		ctx.Next()
	}
}

func ParseTokenFromHeader(cfg *config.Config, ctx *gin.Context) (string, error) {
	var (
		headerToken = ctx.Request.Header.Get("Authorization")
		splitToken  []string
	)

	splitToken = strings.Split(headerToken, "Bearer ")

	// check valid bearer token
	if len(splitToken) <= 1 {
		return "", apperr.NewWithCode(apperr.CodeHTTPUnauthorized, `Invalid Token`)
	}

	return splitToken[1], nil
}
