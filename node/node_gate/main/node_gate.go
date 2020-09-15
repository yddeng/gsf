package main

import (
	"fmt"
	"github.com/yddeng/gsf/cluster"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/node/common/config"
	"github.com/yddeng/gsf/node/node_gate/gate"
	"github.com/yddeng/gsf/util"
	"github.com/yddeng/gsf/util/signal"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Sprintf("args: p, configPath"))
	}

	// 加载配置
	path := os.Args[1]
	conf := util.Must(config.LoadConfig(path)).(*config.Config)

	// 日志
	filename := fmt.Sprintf("gate_%s", conf.Gate.LogicAddr)
	util.InitLogger("log", filename)

	// 集群中启动
	selfAddr := util.Must(addr.MakeAddr(conf.Gate.LogicAddr, conf.Gate.ClusterAddr)).(*addr.Addr)
	util.Must(nil, cluster.Launch(conf.Common.CenterAddr, selfAddr))

	// 启动 gate, 对外服务
	util.Must(nil, gate.Launch(conf.Gate.ExternalAddr))

	util.Logger().Infof("receive signal:%s to shutdown", <-signal.ListenStop())
}
