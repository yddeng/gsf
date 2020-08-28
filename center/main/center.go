package main

import (
	"fmt"
	"github.com/yddeng/gsf/center"
	_ "github.com/yddeng/gsf/codec/cs"
	_ "github.com/yddeng/gsf/codec/ss"
	"github.com/yddeng/gsf/util"
	"github.com/yddeng/gsf/util/signal"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Sprint("program, netAddr"))
		return
	}

	util.InitLogger("log", "center")

	center.Launcher(os.Args[1])
	util.Logger().Infof("receive signal:%s to shutdown", <-signal.ListenStop())
}
