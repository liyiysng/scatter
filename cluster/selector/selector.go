package selector

import (
	"github.com/liyiysng/scatter/logger"
)

// DefaultPolicy 缺省选择策略
const DefaultPolicy = "round_robin"

var (
	myLog = logger.Component("selector")
)
