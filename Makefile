gen_ss:
	cd protocol/ss;make;cd ../../
gen_cs:
	cd protocol/cs;make;cd ../../../
gen_rpc:
	cd protocol/rpc;make;cd ../../../
proto:
	make gen_ss;make gen_cs;make gen_rpc;

