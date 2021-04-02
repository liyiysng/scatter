package nsq_ps

import (
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	nsqd_addr        = "127.0.0.1:4150"
	nsqd_lookup_addr = "127.0.0.1:4161"
)

func TestNSQ(t *testing.T) {
	// producer, err := nsq.NewProducer(nsqd_addr, nsq.NewConfig())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer producer.Stop()
	// err = producer.Publish("topic", []byte("foo 6"))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second * 11 // 每11秒轮询更新一次服务地址

	consumer, err := nsq.NewConsumer("topic", "xxx", cfg)

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		t.Logf("recv %s", string(message.Body))
		return nil
	}))

	err = consumer.ConnectToNSQLookupd(nsqd_lookup_addr)
	if err != nil {
		t.Fatal(err)
	}

	defer consumer.Stop()

	go func() {
		time.Sleep(time.Second * 10)
		consumer.Stop()
	}()

	<-consumer.StopChan
}
