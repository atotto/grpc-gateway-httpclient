module github.com/atotto/grpc-gateway-httpclient/test

go 1.12

require (
	github.com/akutz/memconn v0.1.0
	github.com/atotto/grpc-gateway-httpclient v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f
	google.golang.org/grpc v1.53.0
)

replace github.com/atotto/grpc-gateway-httpclient => ../
