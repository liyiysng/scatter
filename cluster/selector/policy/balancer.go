package policy

import (
	"github.com/liyiysng/scatter/logger"
	"google.golang.org/grpc/balancer"
)

var (
	myLog = logger.Component("policy")
)

func init() {
	balancer.Register(newSessionAffinityBuilder())
}
