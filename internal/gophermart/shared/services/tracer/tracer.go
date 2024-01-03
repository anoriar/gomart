package tracer

import (
	"fmt"
	serverConfig "github.com/anoriar/gophermart/internal/gophermart/shared/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

func NewTracer(serverConfig *serverConfig.Config) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serverConfig.TracerServiceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
		Headers: &jaeger.HeadersConfig{
			TraceContextHeaderName: serverConfig.TracerHeader,
		},
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Metrics(metrics.NullFactory))

	if err != nil {
		return nil, nil, fmt.Errorf("failed get jaeger config from env %w", err)
	}

	return tracer, closer, nil
}
