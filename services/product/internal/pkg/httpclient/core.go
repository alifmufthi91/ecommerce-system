package httpclient

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/observ"
	resty "github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func (c *restyC) BuildRequest(ctx context.Context, prop *PropRequest) (*resty.Request, map[string]string) {
	var (
		q     = c.client.R().SetContext(ctx)
		attrs = make(map[string]string)
	)

	if prop.Headers != nil {
		q.SetHeaders(prop.Headers)
		filtered := []string{"authorization", "cookie"}
		// http header attributes
		for k, h := range prop.Headers {
			key := strings.ToLower(k)
			isFiltered := false
			for _, v := range filtered {
				if strings.EqualFold(key, v) {
					isFiltered = true
					break
				}
			}
			if isFiltered {
				continue
			}
			attrs["request."+key] = h
		}

	}

	if prop.Body != nil {
		q.SetBody(prop.Body)
		attrs["request.body"] = c.PrettyPrint(prop.Body)

	}

	if prop.URIParams != nil {
		q.SetPathParams(prop.URIParams)
		attrs["request.uri.params"] = c.PrettyPrint(prop.URIParams)

	}

	if prop.QueryParams != nil {
		q.SetQueryParams(prop.QueryParams)
		attrs["request.query.params"] = c.PrettyPrint(prop.QueryParams)

	}

	if prop.QueryString != "" {
		q.SetQueryString(prop.QueryString)
		attrs["request.query"] = prop.QueryString

	}

	if prop.FormData != nil {
		q.SetFormData(prop.FormData)
		attrs["request.formdata"] = c.PrettyPrint(prop.FormData)
	}

	if prop.MultiFormData != nil {
		q.SetFormDataFromValues(prop.MultiFormData)
	}

	if prop.MultiQueryParams != nil {
		q.SetQueryParamsFromValues(prop.MultiQueryParams)
	}

	if prop.Files != nil {
		q.SetFiles(prop.Files)
	}

	if prop.FileReaders != nil {
		for _, f := range prop.FileReaders {
			q = q.SetFileReader(f.Param, f.FileName, f.Reader)
		}
	}

	return q, attrs
}

func (c *restyC) do(ctx context.Context, method string, prop *PropRequest) (*http.Response, error) {
	var tracer = c.getTracer(ctx)

	ctx, span := tracer.Start(ctx, `httpclient.do`)
	defer span.End()

	span.SetAttributes(
		attribute.String("http.uri", prop.URI),
		attribute.String("http.method", method),
	)

	attrs := []attribute.KeyValue{}

	defer func() {
		span.AddEvent("HTTP Request Log", trace.WithAttributes(attrs...))
	}()

	if prop.WithRetry {

		if prop.MaxRetry != 0 {
			c.client.SetRetryCount(prop.MaxRetry)
		}

		if prop.RetryWaitTime != 0 {
			retryWaitTime := time.Duration(prop.RetryWaitTime) * time.Second
			c.client.SetRetryWaitTime(retryWaitTime).SetRetryMaxWaitTime(retryWaitTime)
		}

		c.client.AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusInternalServerError
			},
		)
	}

	q, ats := c.BuildRequest(ctx, prop)

	// add tracing context to HTTP header
	observ.InjectTraceHeader(ctx, q.Header)

	for k, attr := range ats {
		attrs = append(attrs, attribute.String(k, attr))
	}

	resp, err := q.Execute(method, prop.URI)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		attrs = append(attrs, attribute.String("response.status", resp.RawResponse.Status))
		dumpresponse, _ := httputil.DumpResponse(resp.RawResponse, true)
		attrs = append(attrs, attribute.String("response.dumpresponse", string(dumpresponse)))
	}

	return resp.RawResponse, err
}

func (c *restyC) doJSON(ctx context.Context, method string, prop *PropRequest) (*http.Response, error) {
	if prop.Headers == nil {
		prop.Headers = make(map[string]string)
	}
	prop.Headers[ContentType] = ContentJSON
	prop.Headers[ContentAccept] = ContentJSON
	return c.do(ctx, method, prop)
}

func (c *restyC) GetJSON(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.doJSON(ctx, http.MethodGet, prop)
}

func (c *restyC) PostJSON(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.doJSON(ctx, http.MethodPost, prop)
}

func (c *restyC) PutJSON(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.doJSON(ctx, http.MethodPut, prop)
}

func (c *restyC) PatchJSON(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.doJSON(ctx, http.MethodPatch, prop)
}

func (c *restyC) Get(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, prop)
}

func (c *restyC) Post(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, prop)
}

func (c *restyC) Put(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, prop)
}

func (c *restyC) Delete(ctx context.Context, prop *PropRequest) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, prop)
}

func (c *restyC) getTracer(ctx context.Context) oteltrace.Tracer {
	var tracer oteltrace.Tracer
	// get global tracer
	tracerInterface := ctx.Value(tracerKey)
	if tracerInterface != nil {
		tc, _ := tracerInterface.(oteltrace.Tracer)
		tracer = tc

	}

	if tracer == nil {
		// if not found make one
		tracer = otel.GetTracerProvider().Tracer(
			tracerName,
		)
	}

	return tracer
}

func (c *restyC) PrettyPrint(data interface{}) string {
	JSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println("error marshal pretty print : ", err)
	}

	return string(JSON)
}
