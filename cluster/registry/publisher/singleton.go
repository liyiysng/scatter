package publisher

import "sync"

var defaultPublisher *Publisher

var initPubOnce sync.Once

// GetPublisher 获取单例
func GetPublisher() *Publisher {

	initPubOnce.Do(func() {
		pub, err := NewPublisher(defaultPubOptions)
		if err != nil {
			myLog.Errorf("init default publish error : %v", err)
			return
		}
		defaultPublisher = pub
	})

	if defaultPublisher == nil {
		panic("init default publish failed")
	}

	return defaultPublisher
}
