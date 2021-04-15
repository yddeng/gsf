package main

import (
	"fmt"
	"github.com/yddeng/clugs/cluster"
	_ "github.com/yddeng/clugs/codec/cs"
	_ "github.com/yddeng/clugs/codec/ss"
	"github.com/yddeng/clugs/logger"
	"github.com/yddeng/clugs/util/signal"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Sprint("program, netAddr"))
		return
	}

	logge := logger.New("log", "center")
	logger.InitLogger(logge)

	center := cluster.LunchCenter(os.Args[1])

	logger.Infof("receive signal:%s to shutdown", <-signal.ListenStop())
	center.Stop()
}
