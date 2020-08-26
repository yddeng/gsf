all:
	test -d ../rpc || mkdir -p ../rpc
	rm ../rpc/*.go; protoc -I ./ -I ../../.. --go_out=paths=source_relative:../rpc *.proto
