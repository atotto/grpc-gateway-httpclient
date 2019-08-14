package service

import (
	"context"

	"github.com/atotto/grpc-gateway-httpclient/test/testdata/apis"
)

var _ apis.EchoServiceServer = (*EchoService)(nil)

type EchoService struct {
}

func NewEchoService() *EchoService {
	return &EchoService{}
}

func (s *EchoService) Echo(ctx context.Context, req *apis.EchoRequest) (res *apis.EchoResponse, err error) {
	return &apis.EchoResponse{Echo: req.GetEcho()}, nil
}
