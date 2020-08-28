all:
	test -d ../rpcpb || mkdir -p ../rpcpb
	rm ../rpcpb/*.go; protoc -I ./ -I ../../.. --go_out=paths=source_relative:../rpcpb *.proto
