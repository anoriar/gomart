package tracer

import (
	"github.com/anoriar/gophermart/internal/gophermart/shared/middleware/tracer/internal/responsewriter"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
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

		tw := responsewriter.NewTracerResponseWriter(w)

		h.ServeHTTP(tw, request.WithContext(ctx))

		responseData := tw.ResponseData()

		span.SetTag("http.status_code", responseData.Status())
		if responseData.Status() >= http.StatusInternalServerError {
			ext.Error.Set(span, true)
		}
		span.LogFields(log.String("payload", string(responseData.Payload())))
	})
}
