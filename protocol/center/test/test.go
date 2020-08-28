package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/protocol/center/center"
)

func main() {
	s := &center.LoginResp{
		Ok: true,
	}

	b, err := proto.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b, err)

	var e center.LoginResp
	err = proto.Unmarshal(b, &e)
	fmt.Println(e, err)

	k := ss()
	fmt.Println(k)
	fmt.Println("ss")

}

func ss() int {
	k := 1
	defer func() {
		k = 2
	}()
	return k
}
