install:
	go install ./protoc-gen-grpc-gateway-httpclient

.PHONY: test
test: install
	make -C test test
