package node

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkChanelWakeup(b *testing.B) {
	b.ResetTimer()

	goCount := int32(100000)

	total := int64(0)

	waitTimes := int32(0)

	// prepare N chanel
	chanels := make([]chan int64, goCount)

	for i := int32(0); i < goCount; i++ {
		chanels[i] = make(chan int64, 100)
	}

	// prepare goCount Go wait for read
	prepareWg := sync.WaitGroup{}
	rwg := sync.WaitGroup{}
	for i := int32(0); i < goCount; i++ {
		g := i
		rwg.Add(1)
		prepareWg.Add(1)
		go func() {
			prepareWg.Done()
			for {
				select {
				case r, ok := <-chanels[g]:
					{
						if !ok {
							goto exit
						}
						atomic.AddInt64(&total, r)
					}
				default:
					{
						atomic.AddInt32(&waitTimes, 1)
						select {
						case r, ok := <-chanels[g]:
							{
								if !ok {
									goto exit
								}
								atomic.AddInt64(&total, r)
							}
						}

					}
				}
			}
		exit:
			rwg.Done()
		}()
	}

	prepareWg.Wait()

	b.Logf("prepare go %d", goCount)

	b.ResetTimer()

	wwg := sync.WaitGroup{}

	// send data
	sendCount := int32(0)
	b.RunParallel(func(tb *testing.PB) {
		for tb.Next() {
			wwg.Add(1)
			id := atomic.AddInt32(&sendCount, 1) % goCount
			r, _ := strconv.ParseInt("1", 10, 64)
			if sendCount%1000 == 0 {
				time.Sleep(time.Microsecond * time.Duration(rand.Int31n(2000)))
			}
			chanels[id] <- r
			wwg.Done()
		}
	})

	wwg.Wait()

	// close all chan
	for i := int32(0); i < goCount; i++ {
		close(chanels[i])
	}

	rwg.Wait()

	b.Logf("send count %d , recv cound %d , wait count %d", sendCount, total, waitTimes)
}

func BenchmarkPlane(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		goCount := int32(100)

		total := int64(0)

		waitTimes := int32(0)

		// prepare N chanel
		chanels := make([]chan int64, goCount)

		for i := int32(0); i < goCount; i++ {
			chanels[i] = make(chan int64, 100)
		}

		// prepare goCount Go wait for read
		prepareWg := sync.WaitGroup{}
		rwg := sync.WaitGroup{}
		for i := int32(0); i < goCount; i++ {
			g := i
			rwg.Add(1)
			prepareWg.Add(1)
			go func() {

				for {
					select {
					case r, ok := <-chanels[g]:
						{
							if !ok {
								goto exit
							}
							atomic.AddInt64(&total, r)
						}
					default:
						{
							prepareWg.Done()
							atomic.AddInt32(&waitTimes, 1)
							select {
							case r, ok := <-chanels[g]:
								{
									if !ok {
										goto exit
									}
									atomic.AddInt64(&total, r)
								}
							}

						}
					}
				}
			exit:
				rwg.Done()
			}()
		}

		prepareWg.Wait()
	}
}

func BenchmarkNoChannel(b *testing.B) {
	b.ResetTimer()

	total := int64(0)

	wg := sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			r, _ := strconv.ParseInt("10", 10, 64)
			atomic.AddInt64(&total, r)
			wg.Done()
		}()
	}
	wg.Wait()
	b.Logf("%d", total)
}
