package tracer

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

type TracerMiddleware struct {
}

func (tm *TracerMiddleware) Trace(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		opName := request.Method + " " + request.URL.Path
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
		span, ctx := opentracing.StartSpanFromContext(request.Context(), opName, ext.RPCServerOption(spanCtx))
		defer span.Finish()

		h.ServeHTTP(w, request.WithContext(ctx))
	})
}
