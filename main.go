package main

import (
	"github.com/liyiysng/scatter/metrics"

	"log"
)

func main() {

	err := metrics.StartListen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
