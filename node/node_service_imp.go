package node

import (
	"context"

	"github.com/liyiysng/scatter/protobuf/node"
)

// ServiceImp 实现
type ServiceImp struct {
	node.UnimplementedNodeServer
	str string
}

// GetNodeState 获取节点状态
func (imp *ServiceImp) GetNodeState(context.Context, *node.NodeIdentity) (*node.NodeState, error) {
	return nil, nil
}
