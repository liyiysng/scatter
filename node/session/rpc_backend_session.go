package session

import "google.golang.org/protobuf/proto"

// rpcBackendSession 实现BackendSession
// 节点间通过rpc同步的backend session
type rpcBackendSession struct {
}

// NewRPCBackendSession 创建一个rpcBackendSession
func NewRPCBackendSession() BackendSession {
	return nil
}

func (s *rpcBackendSession) BindStrContext(key string) proto.Message {
	return nil
}
