package jaeger

import (
	"io"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// Options holds configuration options for Jaeger
type Options struct {
	Disabled    bool
	Probability float64
	ServiceName string
}

// Configure configures a global Jaeger tracer
func Configure(options Options) (io.Closer, error) {
	cfg := config.Configuration{
		Disabled: options.Disabled,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: options.Probability,
		},
	}
	return cfg.InitGlobalTracer(options.ServiceName)
}
