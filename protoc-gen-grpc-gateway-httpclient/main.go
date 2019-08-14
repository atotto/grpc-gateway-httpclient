package main

import (
	"github.com/gogo/protobuf/vanity/command"
	"github.com/atotto/grpc-gateway-httpclient/protoc-gen-grpc-gateway-httpclient/plugin"
)

func main() {
	req := command.Read()
	p := plugin.NewPlugin()
	resp := command.GeneratePlugin(req, p, ".pb.http.go")
	command.Write(resp)
}
