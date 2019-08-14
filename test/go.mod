module github.com/atotto/grpc-gateway-httpclient/test

go 1.12

require (
	github.com/akutz/memconn v0.1.0
	github.com/atotto/grpc-gateway-httpclient v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.9.5
	google.golang.org/genproto v0.0.0-20190801165951-fa694d86fc64
	google.golang.org/grpc v1.23.0
)

replace github.com/atotto/grpc-gateway-httpclient => ../
