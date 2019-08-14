

gen:
	protoc -I testdata/protobuf -I testdata/apis --grpc-gateway-httpclient_out=testdata/apis --go_out=plugins=grpc:testdata/apis --grpc-gateway_out=testdata/apis testdata/apis/echo.proto
	protoc -I testdata/protobuf -I testdata/apis --grpc-gateway-httpclient_out=testdata/apis --go_out=plugins=grpc:testdata/apis --grpc-gateway_out=testdata/apis testdata/apis/message.proto

test: gen
	go test -v
