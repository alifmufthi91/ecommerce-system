package productservice

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/alifmufthi91/ecommerce-system/services/order/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/httpclient"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
)

//go:generate mockery --name=IProductSvc --case underscore
type IProductSvc interface {
	GetProductByID(ctx context.Context, req GetProductByIDReq) (GetProductByIDResp, error)
}

type ProductSvc struct {
	URL        string
	httpClient httpclient.IHTTPClient
}

const (
	ProductServicePrefix = "productservice: "
)

var _ IProductSvc = (*ProductSvc)(nil)

func Init(opts _options.DefaultOptions) *ProductSvc {
	return &ProductSvc{
		URL:        opts.Config.External.ProductServiceURL,
		httpClient: opts.HttpClient,
	}
}

func (w *ProductSvc) GetProductByID(ctx context.Context, req GetProductByIDReq) (res GetProductByIDResp, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "productsvc.GetProductByID")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	resp, err := w.httpClient.Get(ctx, &httpclient.PropRequest{
		URI: w.URL + "/products/" + req.ProductID,
		Headers: map[string]string{
			"Authorization": "Bearer " + req.Token,
		},
	})

	if err != nil {
		return res, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, ProductServicePrefix+`internal server error`)
	}

	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusMultipleChoices {
		if err := handleErrorResponse(rawData, resp.StatusCode); err != nil {
			return res, err
		}
	}

	if err = json.Unmarshal(rawData, &res); err != nil {
		return res, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, ProductServicePrefix+"failed to unmarshal response")
	}

	return res, nil
}

func handleErrorResponse(rawData []byte, statusCode int) error {
	if len(rawData) == 0 {
		return apperr.NewWithCode(apperr.MapStatusCodeToErrorCode(statusCode), ProductServicePrefix+http.StatusText(statusCode))
	}

	var failResp ErrorResponse

	if err := json.Unmarshal(rawData, &failResp); err != nil {
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, ProductServicePrefix+"failed to unmarshal error response")
	}

	return apperr.NewWithCode(apperr.MapStatusCodeToErrorCode(statusCode), ProductServicePrefix+failResp.Metadata.Message)
}
