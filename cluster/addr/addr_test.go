package addr

import (
	"fmt"
	"testing"
)

func TestMakeAddr(t *testing.T) {
	addr, err := MakeAddr("4095.2.1023", "127.0.0.1:8010")
	if nil != err {
		fmt.Println(err)
	} else {
		fmt.Println(addr.Logic.String(), addr.Logic.Group(), addr.Logic.Type(), addr.Logic.Server(), addr.Net)
	}

}
