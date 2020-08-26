all:
	test -d ../center || mkdir -p ../center
	rm ../center/*.go; protoc -I ./ -I ../../.. --go_out=paths=source_relative:../center *.proto
