package metrics

import "github.com/spf13/viper"

// NewCustomMetricsSpec returns a *CustomMetricsSpec by reading config key
func NewCustomMetricsSpec(config *viper.Viper) (*CustomMetricsSpec, error) {
	var spec CustomMetricsSpec

	err := config.UnmarshalKey("scatter.metrics.custom", &spec)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}
