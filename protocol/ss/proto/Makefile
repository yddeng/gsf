all:
	test -d ../sspb || mkdir -p ../sspb
	rm ../sspb/*.go; protoc -I ./ -I ../../.. --go_out=paths=source_relative:../sspb *.proto
