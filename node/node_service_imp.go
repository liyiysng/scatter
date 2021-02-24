package node

import (
	"context"

	"github.com/liyiysng/scatter/node/grpc_proto"
)

// ServiceImp 实现
type ServiceImp struct {
	grpc_proto.UnimplementedNodeServer
	str string
}

// GetNodeState 获取节点状态
func (imp *ServiceImp) GetNodeState(context.Context, *grpc_proto.NodeIdentity) (*grpc_proto.NodeState, error) {
	return nil, nil
}
