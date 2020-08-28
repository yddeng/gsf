all:
	test -d ../cspb || mkdir -p ../cspb
	rm ../cspb/*.go; protoc -I ./ -I ../../.. --go_out=paths=source_relative:../cspb *.proto
