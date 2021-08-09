package elastic

import (
	"context"
	"sync"
	"time"

	"github.com/liyiysng/scatter/util"
	"github.com/olivere/elastic/v7"
)

type ESSink struct {
	mu sync.Mutex

	index string

	bulk []interface{}

	wg util.WaitGroupWrapper

	client *elastic.Client

	numBulk int

	writeTicker *time.Ticker

	onWritten func(v []interface{})
}

// NewEsSink 新建es sink
func NewEsSink(index string, numBulk int, flushInverval time.Duration, onWritten func(v []interface{}), options ...elastic.ClientOptionFunc) (s *ESSink, err error) {

	// index create

	// Create an Elasticsearch client
	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	es := &ESSink{
		index:     index,
		client:    client,
		numBulk:   numBulk,
		bulk:      make([]interface{}, 0, numBulk),
		onWritten: onWritten,
	}

	es.startFlushGoroutine(flushInverval)

	return es, err
}

func (es *ESSink) Write(data interface{}) error {
	es.mu.Lock()
	es.bulk = append(es.bulk, data)
	if len(es.bulk) >= es.numBulk {
		es.writeBulk()
	}
	es.mu.Unlock()
	return nil
}

// hold mu
func (es *ESSink) writeBulk() {
	wBulk := es.bulk
	es.bulk = make([]interface{}, 0, es.numBulk)

	es.wg.Wrap(func() {
		defer func() {
			if es.onWritten != nil {
				es.onWritten(wBulk)
			}
		}()
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

func (es *ESSink) flush() {
	es.mu.Lock()
	if len(es.bulk) > 0 {
		es.writeBulk()
	}
	es.mu.Unlock()
}

func (es *ESSink) startFlushGoroutine(interval time.Duration) {
	es.writeTicker = time.NewTicker(interval)
	go func() {
		for range es.writeTicker.C {
			es.flush()
		}
	}()
}

func (es *ESSink) Close() error {
	if es.writeTicker != nil {
		es.writeTicker.Stop()
	}
	es.flush()
	es.wg.Wait()
	return nil
}
