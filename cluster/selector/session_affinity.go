package selector

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type ContextService struct {
}

// BindContext 绑定上下文
func BindContext(ctx context.Context, key string, data proto.Message) error {
	return nil
}

// GetContext 获取上下文
func GetContext(key string) (exists bool, data proto.Message) {
	return false, nil
}
