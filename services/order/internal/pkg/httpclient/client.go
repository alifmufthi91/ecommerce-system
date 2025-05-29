package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/order/config"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	resty "github.com/go-resty/resty/v2"
)

const (
	defaultTimeout  = 60 * time.Second
	infoRequestDump = `httpclient sent request: client_uri=%s request=%s`

	tracerKey  = "otel-go-contrib-tracer"
	tracerName = "httpclient"
)

//go:generate mockery --name=IHTTPClient --case underscore
type IHTTPClient interface {
	Get(ctx context.Context, prop *PropRequest) (*http.Response, error)
	Post(ctx context.Context, prop *PropRequest) (*http.Response, error)
	Put(ctx context.Context, prop *PropRequest) (*http.Response, error)
	Delete(ctx context.Context, prop *PropRequest) (*http.Response, error)

	GetJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
	PostJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
	PutJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
	PatchJSON(ctx context.Context, prop *PropRequest) (*http.Response, error)
}

type restyC struct {
	client *resty.Client
	logger *pkg.Logger
	config *config.Config
}

type Options struct {
	Config *config.Config
	Logger *pkg.Logger
}

func Init(params Options) IHTTPClient {
	r := &restyC{
		client: resty.New(),
		logger: params.Logger,
		config: params.Config,
	}

	r.initClient()

	return r
}

func (c *restyC) initClient() {
	rl := restyLogger{c.logger}

	timeout := defaultTimeout

	c.client.SetTimeout(timeout).
		SetContentLength(true).
		SetCloseConnection(false).
		SetJSONEscapeHTML(true).
		SetDoNotParseResponse(true).
		OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
			//append header x-request-id
			reqID := req.Context().Value(RequestID)
			if reqID != nil {
				req.SetHeader(RequestID, reqID.(string))
			}

			return nil
		}).
		SetPreRequestHook(func(client *resty.Client, req *http.Request) error {
			//dump request
			if req != nil {
				uri := req.URL

				dumpRequest, err := httputil.DumpRequestOut(req, true)
				if err != nil {
					c.logger.WithContext(req.Context()).Warn(err)
				} else {
					c.logger.WithContext(req.Context()).Info(fmt.Sprintf(infoRequestDump, uri, string(dumpRequest)))
				}
			}

			return nil
		}).
		SetLogger(rl)

}
