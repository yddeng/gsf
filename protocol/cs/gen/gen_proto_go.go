package main

import (
	"fmt"
	"github.com/yddeng/gsf/protocol/cs/proto_def"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var message_template = `syntax = "proto3";
option go_package = "gsf/protocol/cs/cspb";

message %sToS {}

message %sToC {}`

func gen_proto(out_path string) {

	fmt.Printf("gen_proto message ............\n")

	for _, v := range proto_def.CS_message {
		filename := fmt.Sprintf("%s/%s.proto", out_path, v.Name)
		//检查文件是否存在，如果存在跳过不存在创建
		f, err := os.Open(filename)
		if err != nil && os.IsNotExist(err) {
			_ = os.MkdirAll(path.Dir(filename), os.ModePerm)
			err := ioutil.WriteFile(filename, []byte(fmt.Sprintf(message_template, v.Name, v.Name)), os.ModePerm)
			if nil != err {
				fmt.Printf("------ error -------- %s Write error:%s\n", v.Name, err.Error())
			}
		} else if nil != f {
			fmt.Printf("%s.proto exist skip\n", v.Name)
			_ = f.Close()
		}
	}

}

var register_template = `package cs

import (
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/protocol/cs/cspb"
)

const (
	CS_SPACE = "c2s"
	SC_SPACE = "s2c"
) 

const(
%s
)

func init() {
	//toS
%s
	//toC
%s
}
`

//产生协议注册文件
func gen_register(out_path string) {
	cmds := ""
	toS_str := ""
	toC_str := ""

	nameMap := map[string]bool{}
	idMap := map[int]bool{}

	for _, v := range proto_def.CS_message {
		if ok, _ := nameMap[v.Name]; ok {
			panic("duplicate message:" + v.Name)
		}

		if ok, _ := idMap[v.Cmd]; ok {
			panic(fmt.Sprintf("duplicate cmd: %d", v.Cmd))
		}

		nameMap[v.Name] = true
		idMap[v.Cmd] = true

		cmds += fmt.Sprintf(`	%s = %d`, strings.Title(v.Name), v.Cmd) + "\n"
		toS_str = toS_str + fmt.Sprintf(`	pb.Register(CS_SPACE,&cspb.%sToS{},%d)`, strings.Title(v.Name), v.Cmd) + "\n"
		toC_str = toC_str + fmt.Sprintf(`	pb.Register(SC_SPACE,&cspb.%sToC{},%d)`, strings.Title(v.Name), v.Cmd) + "\n"
	}

	content := fmt.Sprintf(register_template, cmds, toS_str, toC_str)

	_ = os.MkdirAll(path.Dir(out_path), os.ModePerm)
	err := ioutil.WriteFile(out_path, []byte(content), os.ModePerm)
	if nil != err {
		fmt.Printf("------ error -------- %s Write error:%s\n", out_path, err.Error())
	} else {
		fmt.Printf("%s Write ok\n", out_path)
	}

}

func main() {
	gen_proto("../proto")
	gen_register("../register.go")
	fmt.Printf("cs gen_proto_go ok!\n")
}
