package selector

// DefaultPolicy 缺省选择策略
const DefaultPolicy = "round_robin"

// Selector 选择一个节点
type Selector interface {
	AddService()
	AddNode()
	DeleteNode()
}

// ISession session 定义
type ISession interface {
	// GetSID 获取session id
	GetSID() int64
	// // 绑定上下文
	// BindContext(ctx context.Context, key string, data proto.Message) error
	// // 获取上下文
	// GetContext(key string) (exists bool, data proto.Message)
}

// SessionAffinitySelector 与session关联的selector
type SessionAffinitySelector interface {
	BeginSession()
	EndSession()
}
