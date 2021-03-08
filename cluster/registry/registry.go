// Package registry is an interface for service discovery
package registry

import (
	"errors"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("registry")
)

var (
	// ErrNotFound Not found error when GetService is called
	ErrNotFound = errors.New("service not found")
	// ErrWatcherStopped Watcher stopped error when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// Registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

// Service 服务
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Node 节点
type Node struct {
	ID       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Endpoint service.method
type Endpoint struct {
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// // Value 描述一个值
// type Value struct {
// 	Name   string   `json:"name"`
// 	Type   string   `json:"type"`
// 	Values []*Value `json:"values,omitempty"`
// }

// Option 初始化选项类型
type Option func(*Options)

// RegisterOption 注册选项类型
type RegisterOption func(*RegisterOptions)

// WatchOption 观察选项类型
type WatchOption func(*WatchOptions)

// DeregisterOption 移除注册选项类型
type DeregisterOption func(*DeregisterOptions)

// GetOption 获取服务选项类型
type GetOption func(*GetOptions)

// ListOption 列举服务选项类型
type ListOption func(*ListOptions)

func addNodes(old, neu []*Node) []*Node {
	nodes := make([]*Node, len(neu))
	// add all new nodes
	for i, n := range neu {
		node := *n
		nodes[i] = &node
	}

	// look at old nodes
	for _, o := range old {
		var exists bool

		// check against new nodes
		for _, n := range nodes {
			// ids match then skip
			if o.ID == n.ID {
				exists = true
				break
			}
		}

		// keep old node
		if !exists {
			node := *o
			nodes = append(nodes, &node)
		}
	}

	return nodes
}

func delNodes(old, del []*Node) []*Node {
	var nodes []*Node
	for _, o := range old {
		var rem bool
		for _, n := range del {
			if o.ID == n.ID {
				rem = true
				break
			}
		}
		if !rem {
			nodes = append(nodes, o)
		}
	}
	return nodes
}

// CopyService make a copy of service
func CopyService(service *Service) *Service {
	// copy service
	s := new(Service)
	*s = *service

	// copy nodes
	nodes := make([]*Node, len(service.Nodes))
	for j, node := range service.Nodes {
		n := new(Node)
		*n = *node
		nodes[j] = n
	}
	s.Nodes = nodes

	// copy endpoints
	eps := make([]*Endpoint, len(service.Endpoints))
	for j, ep := range service.Endpoints {
		e := new(Endpoint)
		*e = *ep
		eps[j] = e
	}
	s.Endpoints = eps
	return s
}

// Copy makes a copy of services
func Copy(current []*Service) []*Service {
	services := make([]*Service, len(current))
	for i, service := range current {
		services[i] = CopyService(service)
	}
	return services
}

// Merge merges two lists of services and returns a new copy
func Merge(olist []*Service, nlist []*Service) []*Service {
	var srv []*Service

	for _, n := range nlist {
		var seen bool
		for _, o := range olist {
			if o.Version == n.Version {
				sp := new(Service)
				// make copy
				*sp = *o
				// set nodes
				sp.Nodes = addNodes(o.Nodes, n.Nodes)

				// mark as seen
				seen = true
				srv = append(srv, sp)
				break
			} else {
				sp := new(Service)
				// make copy
				*sp = *o
				srv = append(srv, sp)
			}
		}
		if !seen {
			srv = append(srv, Copy([]*Service{n})...)
		}
	}
	return srv
}

// Remove removes services and returns a new copy
func Remove(old, del []*Service) []*Service {
	var services []*Service

	for _, o := range old {
		srv := new(Service)
		*srv = *o

		var rem bool

		for _, s := range del {
			if srv.Version == s.Version {
				srv.Nodes = delNodes(srv.Nodes, s.Nodes)

				if len(srv.Nodes) == 0 {
					rem = true
				}
			}
		}

		if !rem {
			services = append(services, srv)
		}
	}

	return services
}

//CreateFun 创建一个服务注册器
type CreateFun func(opts ...Option) Registry

var registeredRegistry = make(map[string]CreateFun)

// RegisterRegistry 注册一个服务注册器
func RegisterRegistry(name string, create CreateFun) {
	registeredRegistry[name] = create
}

// GetRegistry 更具名称获取服务注册
func GetRegistry(name string) CreateFun {
	return registeredRegistry[name]
}
