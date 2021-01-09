package main

import (
	"math/rand"
	"time"

	"github.com/liyiysng/scatter/metrics"

	"log"
)

func main() {

	go func() {

		method := []string{"foo", "bar"}
		endpoint := []string{"endpoint:8080", "endpoint:9090"}

		for {
			m := method[rand.Int31n(int32(len(method)))]
			e := method[rand.Int31n(int32(len(endpoint)))]
			log.Printf("sent %s %s\n", m, e)
			metrics.UserInput(m, e)
			time.Sleep(time.Millisecond * time.Duration(rand.Int31n(2000)))
		}
	}()

	err := metrics.StartListen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
