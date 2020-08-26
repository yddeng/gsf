package proto_def

type st struct {
	Name string
	Desc string
	Cmd  int
}

var CS_message = []st{
	st{"heartbeat", "心跳", 1},
	st{"echo", "测试用回射协议", 2},
}
