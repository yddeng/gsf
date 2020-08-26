all:
	test -d ../ss || mkdir -p ../ss
	rm ../ss/*.go; protoc -I ./ -I ../../.. --go_out=paths=source_relative:../ss *.proto
