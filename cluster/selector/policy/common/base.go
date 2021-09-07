package common

import (
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

var (
	myLog = logger.Component("policy/common")
)

// PickerBuilder creates balancer.Picker.
type PickerBuilder interface {
	// Build returns a picker that will be used by gRPC to pick a SubConn.
	Build(info PickerBuildInfo) balancer.Picker
}

// PickerBuildInfo contains information needed by the picker builder to
// construct a picker.
type PickerBuildInfo struct {
	// ReadySCs is a map from all ready SubConns to the Addresses used to
	// create them.
	ReadySCs map[balancer.SubConn]SubConnInfo
	Target   resolver.Target
}

// SubConnInfo contains information about a SubConn created by the base
// balancer.
type SubConnInfo struct {
	Address resolver.Address // the address used to create this SubConn
}

// Config contains the config info about the base balancer builder.
type Config struct {
	// HealthCheck indicates whether health checking should be enabled for this specific balancer.
	HealthCheck bool
}

// NewBalancerBuilder returns a base balancer builder configured by the provided config.
func NewBalancerBuilder(name string, pb PickerBuilder, config Config) balancer.Builder {
	return &baseBuilder{
		name:          name,
		pickerBuilder: pb,
		config:        config,
	}
}

// NewBalancerBuilderWithConnectivityStateEvaluator returns a base balancer builder configured by the provided config.
func NewBalancerBuilderWithConnectivityStateEvaluator(name string, pb PickerBuilder, config Config, stateEvaluator IConnectivityStateEvaluator) balancer.Builder {
	return &baseBuilder{
		name:          name,
		pickerBuilder: pb,
		config:        config,
		csEvltr:       stateEvaluator,
	}
}
