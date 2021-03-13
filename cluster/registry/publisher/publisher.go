package publisher

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/util"
)

var (
	myLog = logger.Component("publisher")
)

// PubOptions 发布选项
type PubOptions struct {
	timeout    time.Duration
	bufferSize int
	// 检测过滤,过滤无需关注的服务,或者节点
	// 默认跳过本节点的服务
	watchFillter func(srvName string, n *registry.Node) bool
}

var defaultPubOptions = &PubOptions{
	timeout:      time.Second * 30,
	bufferSize:   1024,
	watchFillter: func(srvName string, n *registry.Node) bool { return !strings.Contains(srvName, "consul") }, // 过滤consul服务
}

type pubData struct {
	srvName string
	node    *registry.Node
}

// Publisher 发布注册事件
type Publisher struct {
	opts    *PubOptions
	pub     *util.Publisher
	watcher registry.Watcher

	srvs map[string] /*name*/ *registry.Service
	mu   sync.Mutex
}

// NewPublisher 创建发布
func NewPublisher(ops *PubOptions) (pub *Publisher, err error) {

	watcher, err := registry.Watch()
	if err != nil {
		return nil, err
	}

	ret := &Publisher{
		opts:    ops,
		pub:     util.NewPublisher(ops.timeout, ops.bufferSize),
		srvs:    map[string]*registry.Service{},
		watcher: watcher,
	}

	go ret.watchSrvs(watcher)

	return ret, nil
}

// Subscribe 订阅
// 复制写入当前已改变的服务状态
func (s *Publisher) Subscribe(filtter func(srvName string, node *registry.Node) bool) chan interface{} {
	return s.pub.SubscribeTopic(func(v interface{}) bool {
		data := v.(*pubData)
		return filtter(data.srvName, data.node)
	})
}

// Evict 取消订阅
func (s *Publisher) Evict(sub chan interface{}) {
	s.pub.Evict(sub)
}

// FindAllNodes 获取所有满足条件的节点
func (s *Publisher) FindAllNodes(match func(srv *registry.Service, node *registry.Node) bool) []*registry.Node {
	ret := []*registry.Node{}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sr := range s.srvs {
		for _, n := range sr.Nodes {
			if match(sr, n) {
				ret = append(ret, n)
			}
		}
	}

	return ret
}

// Close 关闭
func (s *Publisher) Close() {
	s.pub.Close()
	s.watcher.Stop()
}

func (s *Publisher) watchSrvs(watcher registry.Watcher) {
	for {
		res, err := watcher.Next()
		if err == registry.ErrWatcherStopped {
			return
		}

		if len(res.Service.Nodes) == 0 { // 若无节点数据,无需关心
			continue
		}

		// filltter
		nodes := make([]*registry.Node, 0, len(res.Service.Nodes)) // copy remove
		for _, n := range res.Service.Nodes {
			if !s.opts.watchFillter(res.Service.Name, n) {
				continue
			}
			nodes = append(nodes, n)
		}

		if len(nodes) == 0 {
			continue
		}

		switch res.Action {
		case "create":
			//1.{"Action":"create","Service":{"name":"SrvStrings","version":"","endpoints":null,"nodes":null}}
			//2.{"Action":"create","Service":{"name":"SrvStrings","version":"0.0.1","endpoints":[{"name":"ToLower"},{"name":"ToUpper"},{"name":"Split"}],"nodes":null}}
			//3.{"Action":"update","Service":{"name":"SrvStrings","version":"0.0.1","endpoints":[{"name":"ToLower"},{"name":"ToUpper"},{"name":"Split"}],"nodes":[{"id":"11010","address":"127.0.0.1:1155"}]}}
			fallthrough
		case "update": // create or update
			{
				s.onUpdate(res.Service, nodes)
			}
		case "delete":
			{
				s.onDelete(res.Service, nodes)
			}
		default:
			{
				myLog.Errorf("invalid watch action %s", res.Action)
			}
		}
	}
}

func (s *Publisher) onUpdate(newSrv *registry.Service, nodes []*registry.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 比对现有服务
	if srv, ok := s.srvs[newSrv.Name]; ok { // 服务存在
		for _, newNode := range nodes { // 是否有新节点添加
			found := false
			for _, oldNode := range srv.Nodes {
				if oldNode.SrvNodeID == newNode.SrvNodeID { // 已存在改节点 , 做更新操作 , 替换元数据
					oldNode.Metadata = newNode.Metadata
					found = true
					break
				}
			}
			if !found { // 有新节点加入
				s.addNode(newSrv, newNode)
			}
		}

	} else { //创建服务,并且添加节点
		s.addSrv(newSrv)
		for _, n := range nodes {
			s.addNode(newSrv, n)
		}
	}
}

func (s *Publisher) onDelete(newSrv *registry.Service, nodes []*registry.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range nodes {
		s.delNode(newSrv, n)
	}
}

func (s *Publisher) addSrv(srv *registry.Service) {

	gs := &registry.Service{
		Name:      srv.Name,
		Version:   srv.Version,
		Endpoints: srv.Endpoints,
		Metadata:  srv.Metadata,
	}

	s.srvs[srv.Name] = gs
	buf, _ := json.Marshal(gs)
	myLog.Infof("find servcie %s", string(buf))
}

func (s *Publisher) addNode(srv *registry.Service, node *registry.Node) {
	if rs, ok := s.srvs[srv.Name]; ok {
		if len(srv.Endpoints) > 0 {
			rs.Endpoints = srv.Endpoints
		}
		if len(srv.Metadata) > 0 {
			rs.Metadata = srv.Metadata
		}
		myLog.Infof("find node servcice %s : %v", rs.Name, node)
		rs.Nodes = append(rs.Nodes, node)

		s.pub.Publish(&pubData{
			srvName: srv.Name,
			node:    node,
		})
	} else {
		myLog.Errorf("[Publisher.addNode] servers %s not found", srv.Name)
	}
}

func (s *Publisher) delNode(srv *registry.Service, node *registry.Node) {

	if myLog.V(logger.VDEBUG) {
		myLog.Infof("[Publisher.delNode] %v", node)
	}

	if rs, ok := s.srvs[srv.Name]; ok {
		found := false
		for index, n := range rs.Nodes {
			if n.SrvNodeID == node.SrvNodeID {
				found = true
				rs.Nodes = append(rs.Nodes[:index], rs.Nodes[index+1:]...)
				s.pub.Publish(&pubData{
					srvName: srv.Name,
					node:    node,
				})
				break
			}
		}

		if !found {
			myLog.Errorf("[Publisher.delNode] servers %s node %v not found", srv.Name, node)
		}

	} else {
		myLog.Errorf("[Publisher.delNode] servers %s not found", srv.Name)
	}
}
