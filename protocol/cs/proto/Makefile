all:
	test -d ../cs || mkdir -p ../cs
	rm ../cs/*.go; protoc -I ./ -I ../../.. --go_out=paths=source_relative:../cs *.proto
