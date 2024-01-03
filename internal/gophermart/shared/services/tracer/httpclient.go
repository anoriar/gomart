package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

func NewTracerHTTPClient() *http.Client {
	return &http.Client{
		Transport: TracerRoundTripper{Proxy: http.DefaultTransport},
	}
}

type TracerRoundTripper struct {
	Proxy http.RoundTripper
}

func (t TracerRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	opName := request.Method + " " + request.URL.Path

	span, _ := opentracing.StartSpanFromContext(request.Context(), opName)
	defer span.Finish()

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, request.URL.Path)
	ext.HTTPMethod.Set(span, request.Method)
	err := span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header),
	)
	if err != nil {
		return nil, err
	}

	response, err := t.Proxy.RoundTrip(request)
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}

	return response, nil
}
