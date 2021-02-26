package metrics

import "github.com/liyiysng/scatter/config"

// NewCustomMetricsSpec returns a *CustomMetricsSpec by reading config key
func NewCustomMetricsSpec(config *config.Config) (*CustomMetricsSpec, error) {
	var spec CustomMetricsSpec

	err := config.UnmarshalKey("scatter.metrics.custom", &spec)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}
