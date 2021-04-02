package cluster

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/liyiysng/scatter/cluster/registry"
	"github.com/liyiysng/scatter/cluster/registry/publisher"
	"github.com/liyiysng/scatter/cluster/selector"
	"github.com/liyiysng/scatter/cluster/selector/policy"
	"github.com/liyiysng/scatter/cluster/sessionpb"
	"github.com/liyiysng/scatter/cluster/subsrv"
	"github.com/liyiysng/scatter/cluster/subsrvpb"
	"github.com/liyiysng/scatter/encoding"

	//proto编码
	_ "github.com/liyiysng/scatter/encoding/proto"
	"github.com/liyiysng/scatter/logger"
	"github.com/liyiysng/scatter/util"
)

var (
	// ErrorPublisherClosed 已关闭
	ErrorPublisherClosed = errors.New("publisher closed")
)

type PublishOption struct {
	// 日志
	logerr logger.Logger
}

type OptFunType func(*PublishOption)

func OptPublishWithLog(l logger.Logger) OptFunType {
	return func(po *PublishOption) {
		po.logerr = l
	}
}

type PubError struct {
	errs []error
}

func (e *PubError) Error() string {
	return fmt.Sprintf("%v", e.errs)
}

// Publisher 发布
type Publisher struct {
	opt       PublishOption
	clitents  map[string] /*topic*/ subsrvpb.SubServiceClient
	clientsMu sync.Mutex
	codec     encoding.Codec

	clientBuilder *GrpcClient
	closeEvent    *util.Event

	topicNodes map[string][]string
	topicMu    sync.Mutex

	wg util.WaitGroupWrapper
}

func NewPublisher(o ...OptFunType) *Publisher {

	opt := PublishOption{}

	for _, v := range o {
		v(&opt)
	}

	if opt.logerr == nil {
		opt.logerr = myLog
	}

	ret := &Publisher{
		opt:           opt,
		clitents:      make(map[string] /*topic*/ subsrvpb.SubServiceClient),
		clientBuilder: NewGrpcClient(OptGrpcClientWithLogger(opt.logerr)),
		closeEvent:    util.NewEvent(),
		topicNodes:    make(map[string][]string),
		codec:         encoding.GetCodec("proto"),
	}
	ret.wg.Wrap(ret.run, opt.logerr.Errorf)
	return ret
}

// PublishMulti 发布多个
func (p *Publisher) PublishMulti(ctx context.Context, topic string, cmd string, v []interface{}) error {

	if p.closeEvent.HasFired() {
		return ErrorPublisherClosed
	}

	p.topicMu.Lock()
	allNodes, tok := p.topicNodes[topic]
	p.topicMu.Unlock()
	if !tok || len(allNodes) == 0 { // 无订阅者
		return nil
	}

	return nil
}

// Publish 发布
func (p *Publisher) Publish(ctx context.Context, topic string, cmd string, v interface{}) error {

	if p.closeEvent.HasFired() {
		return ErrorPublisherClosed
	}

	// 是否已经是bytes
	var err error
	var data []byte
	if bytesData, ok := v.([]byte); ok {
		data = bytesData
	} else {
		data, err = p.codec.Marshal(v)
		if err != nil {
			return err
		}
	}

	p.topicMu.Lock()
	allNodes, tok := p.topicNodes[topic]
	p.topicMu.Unlock()
	if !tok || len(allNodes) == 0 { // 无订阅者
		return nil
	}

	conn, err := p.getConn(topic)
	if err != nil {
		return err
	}

	perr := &PubError{}
	perrMu := &sync.Mutex{}

	wg := util.WaitGroupWrapper{}

	for _, v := range allNodes {
		ctx = policy.WithNodeID(ctx, v)
		wg.Wrap(func() {
			sinfo := &sessionpb.SessionInfo{
				SType:   sessionpb.SessionType_Pub,
				PubInfo: &sessionpb.PubInfo{},
			}
			res, underlyingError := conn.Pub(ctx, &subsrvpb.PubReq{
				Sinfo:   sinfo,
				Topic:   topic,
				Cmd:     cmd,
				Payload: data,
			})
			perrMu.Lock()
			defer perrMu.Unlock()
			if underlyingError != nil {
				perr.errs = append(perr.errs, underlyingError)
			} else {
				if res.ErrInfo != nil {
					perr.errs = append(perr.errs, errors.New(res.ErrInfo.Err))
				}
			}
		}, p.opt.logerr.Errorf)
	}

	wg.Wait()

	if len(perr.errs) > 0 {
		return perr
	}

	return nil
}

// Close 关闭
func (p *Publisher) Close() {
	p.closeEvent.Fire()
	p.wg.Wait()
}

func (p *Publisher) getConn(topic string) (subsrvpb.SubServiceClient, error) {
	p.clientsMu.Lock()
	defer p.clientsMu.Unlock()
	conn, ok := p.clitents[topic]
	if !ok {
		c, err := p.clientBuilder.GetSubSrvClient(topic, GetClientOptWithPolicy("pub"))
		if err != nil {
			return nil, err
		}
		p.clitents[topic] = c
		conn = c
	}
	return conn, nil
}

func (p *Publisher) run() {
	srvChan := publisher.GetPublisher().Subscribe(func(srvName string, node *registry.Node) bool {
		return srvName == subsrv.SubSrvGrpcName
	})

	defer publisher.GetPublisher().Evict(srvChan)
	for {
		select {
		case <-p.closeEvent.Done():
			{
				p.clientBuilder.Close()
				return
			}
		case changeData, ok := <-srvChan:
			{
				if !ok {
					p.opt.logerr.Error("[Publisher.run] subcrible closed")
					return
				}

				changedSubscribes, err := subsrv.GetSubSrvFromMeta((changeData.(*publisher.PubData)).Node.Metadata)
				if err != nil {
					p.opt.logerr.Errorf("[Publisher.run] %v", err)
					continue
				}

				if len(changedSubscribes) == 0 {
					continue
				}

				for _, changed := range changedSubscribes {
					nodes := publisher.GetPublisher().FindAllNodes(func(srv *registry.Service, node *registry.Node) bool {
						if srv.Name == subsrv.SubSrvGrpcName { // 子服务
							hsrvs, err := subsrv.GetSubSrvFromMeta(node.Metadata)
							if err != nil {
								p.opt.logerr.Errorf("[Publisher.run] node meta error %v", err)
								return false
							}
							found := false
							for _, h := range hsrvs {
								if h == changed {
									found = true
									break
								}
							}
							return found
						}
						return false
					})

					// update
					p.topicMu.Lock()
					if len(nodes) == 0 { // 该主题的所有订阅者都已关闭
						p.topicNodes[changed] = []string{}
					} else {
						nodeIDs := []string{}
						for _, n := range nodes {
							nid, _, err := selector.GetServiceNameAndNodeID(n.SrvNodeID)
							if err != nil {
								p.opt.logerr.Error(err)
								continue
							}
							nodeIDs = append(nodeIDs, nid)
						}
						p.topicNodes[changed] = nodeIDs
					}
					p.topicMu.Unlock()
				}

			}
		}
	}
}
