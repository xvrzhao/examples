package utils

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"log"
)

func InitTracer(serviceName, localAgentHost string, localAgentPort int) io.Closer {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler:     &config.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
		Reporter:    &config.ReporterConfig{LocalAgentHostPort: fmt.Sprintf("%s:%d", localAgentHost, localAgentPort)},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		err = fmt.Errorf("initTracer: newTracer: %w", err)
		log.Fatal(err)
	}

	opentracing.SetGlobalTracer(tracer)
	return closer
}
