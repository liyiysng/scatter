package statsd

import (
	"strings"
)

// HostKey 状态守护进程key的生成
func HostKey(h string) string {
	return strings.Replace(strings.Replace(h, ".", "_", -1), ":", "_", -1)
}
