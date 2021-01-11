package node

// Client 表示一个客户端,可能是TCP/UDP/ws(s)/http(s)/
type Client interface {
	// Stats 获取客户端状态
	Stats() ClientStats
}
