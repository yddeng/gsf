package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/cluster"
	"github.com/yddeng/gsf/cluster/addr"
	protoss "github.com/yddeng/gsf/protocol/ss"
)

func main() {
	cluster.RegisterSSMethod(protoss.Echo, func(from addr.LogicAddr, msg proto.Message) {

	})
}
