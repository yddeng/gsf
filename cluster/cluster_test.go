package cluster

import (
	"fmt"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/util"
	"testing"
)

func TestLauncher(t *testing.T) {
	util.InitLogger("log", "test")
	logic, err := addr.MakeAddr("1.1.1", "127.0.0.1:6547")
	if err != nil {
		fmt.Println(1, err)
		return
	}

	err = Launcher("127.0.0.1:9874", logic)
	if err != nil {
		fmt.Println(2, err)
		return
	}
	select {}
}
