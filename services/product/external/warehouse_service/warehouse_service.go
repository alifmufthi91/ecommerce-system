package warehouseservice

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/alifmufthi91/ecommerce-system/services/product/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/httpclient"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
)

//go:generate mockery --name=IWarehouseSvc --case underscore
type IWarehouseSvc interface {
	GetStockAvailables(ctx context.Context, req GetStockAvailablesReq) (res GetStockAvailablesResp, err error)
}

type WarehouseSvc struct {
	URL        string
	httpClient httpclient.IHTTPClient
}

const (
	WarehouseServicePrefix = "warehouseservice: "
)

var _ IWarehouseSvc = (*WarehouseSvc)(nil)

func Init(opts _options.DefaultOptions) *WarehouseSvc {
	return &WarehouseSvc{
		URL:        opts.Config.External.WarehouseServiceURL,
		httpClient: opts.HttpClient,
	}
}

func (w *WarehouseSvc) GetStockAvailables(ctx context.Context, req GetStockAvailablesReq) (res GetStockAvailablesResp, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "warehousesvc.GetStockAvailables")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	resp, err := w.httpClient.Get(ctx, &httpclient.PropRequest{
		URI:              w.URL + "/stocks/availables",
		MultiQueryParams: req.ToQueryParams(),
		Headers: map[string]string{
			"Authorization": "Bearer " + req.Token,
		},
	})

	if err != nil {
		return res, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, WarehouseServicePrefix+`internal server error`)
	}

	defer resp.Body.Close()

	rawData, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= http.StatusMultipleChoices {
		if err := handleErrorResponse(rawData, resp.StatusCode); err != nil {
			return res, err
		}
	}

	if err = json.Unmarshal(rawData, &res); err != nil {
		return res, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, WarehouseServicePrefix+"failed to unmarshal response")
	}

	return res, nil
}

func handleErrorResponse(rawData []byte, statusCode int) error {
	if len(rawData) == 0 {
		return apperr.NewWithCode(apperr.MapStatusCodeToErrorCode(statusCode), WarehouseServicePrefix+http.StatusText(statusCode))
	}

	var failResp ErrorResponse

	if err := json.Unmarshal(rawData, &failResp); err != nil {
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, WarehouseServicePrefix+"failed to unmarshal error response")
	}

	return apperr.NewWithCode(apperr.MapStatusCodeToErrorCode(statusCode), WarehouseServicePrefix+failResp.Metadata.Message)
}
