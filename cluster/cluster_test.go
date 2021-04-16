package cluster

import (
	"fmt"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/clugs/logger"
	"testing"
)

func init() {
	logge := logger.New("log", "cluster_test")
	logger.InitLogger(logge)
}

func TestLauncher(t *testing.T) {
	logic, err := addr.MakeAddr("1.1.1", "127.0.0.1:6547")
	if err != nil {
		fmt.Println(1, err)
		return
	}

	Launch("127.0.0.1:9874", logic)
	select {}
}
