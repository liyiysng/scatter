package registry

import (
	"fmt"
	"sync"

	"github.com/liyiysng/scatter/config"
)

var defaultRegistryMu sync.Mutex
var defaultRegistry Registry

var initRegistryOnce sync.Once

func getRegistry() Registry {

	initRegistryOnce.Do(func() {

		defaultName := config.GetConfig().GetString("scatter.register.default")

		registryCreater := GetRegistry(defaultName)
		if registryCreater == nil {
			myLog.Errorf("registry scatter.register.default %s not found", defaultName)
			return
		}

		defaultRegistry = registryCreater(
			Addrs(config.GetConfig().GetStringSlice("scatter.register."+defaultName+".addrs")...),
			Timeout(config.GetConfig().GetDuration("scatter.register."+defaultName+".timeout")),
		)
	})

	if defaultRegistry == nil {
		panic(fmt.Errorf("%s initial failed", config.GetConfig().GetString("scatter.register.default")))
	}

	return defaultRegistry
}

// Init 初始化
func Init(opts ...Option) error {
	defaultRegistryMu.Lock()
	defer defaultRegistryMu.Unlock()

	return getRegistry().Init(opts...)
}

// Register 注册
func Register(srv *Service, opts ...RegisterOption) error {
	defaultRegistryMu.Lock()
	defer defaultRegistryMu.Unlock()

	return getRegistry().Register(srv, opts...)
}

// Deregister 取消注册
func Deregister(srv *Service, opts ...DeregisterOption) error {
	defaultRegistryMu.Lock()
	defer defaultRegistryMu.Unlock()

	return getRegistry().Deregister(srv, opts...)
}

// GetService 获取服务
func GetService(srvName string, opts ...GetOption) ([]*Service, error) {
	defaultRegistryMu.Lock()
	defer defaultRegistryMu.Unlock()

	return getRegistry().GetService(srvName, opts...)
}

// ListServices 列出服务
func ListServices(opts ...ListOption) ([]*Service, error) {
	defaultRegistryMu.Lock()
	defer defaultRegistryMu.Unlock()

	return getRegistry().ListServices(opts...)
}

// Watch 监视服务
func Watch(opts ...WatchOption) (Watcher, error) {
	defaultRegistryMu.Lock()
	defer defaultRegistryMu.Unlock()

	return getRegistry().Watch(opts...)
}

// String 名称
func String() string {
	defaultRegistryMu.Lock()
	defer defaultRegistryMu.Unlock()

	return defaultRegistry.String()
}
