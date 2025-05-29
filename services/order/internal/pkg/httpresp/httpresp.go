package httpresp

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/observ"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Meta struct {
	Path       string      `json:"path"`
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Error      error       `json:"error,omitempty" swaggerignore:"true"`
	Timestamp  string      `json:"timestamp"`
	TraceID    string      `json:"trace_id"`
	Data       interface{} `json:"data"`
}

type Response struct {
	Data       interface{} `json:"data"`
	Meta       interface{} `json:"meta,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Success    string      `json:"success"`
	TraceID    string      `json:"trace_id"`
}

type Pagination struct {
	CurrentPage     int64   `json:"current_page"`
	CurrentElements int64   `json:"current_elements"`
	TotalPages      int64   `json:"total_pages"`
	TotalElements   int64   `json:"total_elements"`
	SortBy          string  `json:"sort_by"`
	CursorStart     *string `json:"cursor_start,omitempty"`
	CursorEnd       *string `json:"cursor_end,omitempty"`
}

// HTTPErrResp http error response
type HTTPErrResp struct {
	Meta Meta `json:"metadata"`
}

type HTTPErrRespTest struct {
	Metadata struct {
		Path       string          `json:"path"`
		StatusCode int             `json:"statusCode"`
		Status     string          `json:"status"`
		Message    string          `json:"message"`
		Error      apperr.AppError `json:"error"`
		Timestamp  time.Time       `json:"timestamp"`
		TraceID    string          `json:"trace_id"`
		Data       any             `json:"data"`
	} `json:"metadata"`
}

func HttpRespError(c *gin.Context, err error, data ...interface{}) {
	_, span := observ.GetTracer().Start(c.Request.Context(), "httpresp.HttpRespError")
	defer span.End()

	traceID, spanID := observ.ReadTraceID(c.Request.Context())

	debugMode := false
	lang := "ID"
	if c.Request.Header.Get("x-app-debug") == "true" {
		debugMode = true
	}
	if c.Request.Header.Get("x-app-lang") == "EN" {
		lang = "EN"
	}
	statusCode, displayError := apperr.CompileError(err, lang, debugMode)

	jsonErrResp := &HTTPErrResp{
		Meta: Meta{
			Path:       c.Request.URL.Path,
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
			Message:    fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)),
			Error:      displayError,
			Timestamp:  time.Now().Format(time.RFC3339),
			TraceID:    traceID,
		},
	}

	if len(data) == 1 {
		jsonErrResp.Meta.Data = data[0]
	}

	c.Set("status_code", statusCode)
	c.Set("status", http.StatusText(statusCode))
	c.Set("error", fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)))
	c.Set("trace.id", traceID)
	c.Set("span.id", spanID)

	// ignore default log because its already defined in main.go

	log := pkg.Logger{
		SugaredLogger: zap.L().Sugar(),
	}

	if statusCode == http.StatusUnauthorized {
		log.WithContext(c.Request.Context()).Info(displayError)
	} else if statusCode >= 500 {
		log.WithContext(c.Request.Context()).Error(displayError)
	} else {
		log.WithContext(c.Request.Context()).Warn(displayError)
	}

	c.AbortWithStatusJSON(statusCode, jsonErrResp)
}

func HttpRespSuccess(c *gin.Context, data interface{}, pagination *Pagination) {
	traceID, _ := observ.ReadTraceID(c.Request.Context())

	kind := reflect.ValueOf(data).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		// check kalo data nya nil / kosong
		if data == nil || reflect.ValueOf(data).IsNil() {
			// kalo data arraynya kosong returnnya "data": []
			data = []interface{}{}
		}
	}

	c.JSON(http.StatusOK, Response{
		Data:       data,
		Pagination: pagination,
		Success:    "success",
		TraceID:    traceID,
	})
}
