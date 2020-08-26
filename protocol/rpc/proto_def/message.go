package proto_def

type st struct {
	Name string
	Cmd  int
}

var RPC_message = []st{
	st{"echo", 1},
}
