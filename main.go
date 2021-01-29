package main

import (
	"github.com/liyiysng/scatter/logger"
)

func main() {
	logger.ErrorDepth(2, "wode")
	logger.ErrorDepth(2, "wode")
	logger.ErrorDepth(2, "wode")
	logger.GLogger.Errorln("hello scatter")
}
