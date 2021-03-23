package selector

import (
	"fmt"
	"strings"

	"github.com/liyiysng/scatter/logger"
)

// DefaultPolicy 缺省选择策略
const DefaultPolicy = "round_robin"

var (
	myLog = logger.Component("selector")
)

// GetServiceID 获取serviceID
func GetServiceID(nodeID string, srvName string) string {
	return fmt.Sprintf("%s/%s", nodeID, srvName)
}

// GetServiceNameAndNodeID 根据 service ID 获取service name 和 nodeID
func GetServiceNameAndNodeID(srvID string) (nodeID string, srvName string, err error) {
	srvAndID := strings.Split(srvID, "/")
	if len(srvAndID) != 2 {
		return "", "", fmt.Errorf("invalid format service id %s", srvID)
	}
	return srvAndID[0], srvAndID[1], nil
}

// 属性key,用于从属性中获取相关联的值
const (
	// AttrKeyNodeID 获取节点ID
	AttrKeyNodeID = "nodeID"
	// AttrKeyServiceID 获取服务ID fmt=nodeID/srvName
	AttrKeyServiceID = "serviceID"
	// AttrKeyServiceName 获取服务名
	AttrKeyServiceName = "srvName"
	// AttrKeyMeta 获取meta
	AttrKeyMeta = "meta"
)
