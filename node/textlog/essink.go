package textlog

import (
	"context"
	"sync"
	"time"

	"github.com/liyiysng/scatter/node/session"
	"github.com/liyiysng/scatter/util"
	"github.com/olivere/elastic/v7"
)

const indexPartten = "<scatter-{now/d}-1>"

type esSink struct {
	mu sync.Mutex

	index string

	bulk []*session.MsgInfo

	wg util.WaitGroupWrapper

	client *elastic.Client

	numBulk int

	writeTicker *time.Ticker
}

// NewEsSink 新建es sink
func NewEsSink(index string, numBulk int, flushInverval time.Duration, options ...elastic.ClientOptionFunc) (s Sink, err error) {

	// index create

	// Create an Elasticsearch client
	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	es := &esSink{
		index:   index,
		client:  client,
		numBulk: numBulk,
		bulk:    make([]*session.MsgInfo, 0, numBulk),
	}

	es.startFlushGoroutine(flushInverval)

	return es, err
}

func (es *esSink) Write(info *session.MsgInfo) error {
	es.mu.Lock()
	es.bulk = append(es.bulk, info)
	if len(es.bulk) >= es.numBulk {
		es.writeBulk()
	}
	es.mu.Unlock()
	return nil
}

// hold mu
func (es *esSink) writeBulk() {
	wBulk := es.bulk
	es.bulk = make([]*session.MsgInfo, 0, es.numBulk)

	es.wg.Wrap(func() {
		bulk := es.client.Bulk().Index(es.index)
		for _, v := range wBulk {
			bulk.Add(elastic.NewBulkIndexRequest().Doc(v))
		}
		// commit
		res, err := bulk.Do(context.Background())
		if err != nil {
			myLog.Errorf("write bulk error %v", err)
			return
		}
		if res.Errors {
			myLog.Error("write bulk error")
			return
		}
	}, nil)

}

func (es *esSink) flush() {
	es.mu.Lock()
	if len(es.bulk) > 0 {
		es.writeBulk()
	}
	es.mu.Unlock()
}

func (es *esSink) startFlushGoroutine(interval time.Duration) {
	es.writeTicker = time.NewTicker(interval)
	go func() {
		for range es.writeTicker.C {
			es.flush()
		}
	}()
}

func (es *esSink) Close() error {
	if es.writeTicker != nil {
		es.writeTicker.Stop()
	}
	es.flush()
	es.wg.Done()
	return nil
}
