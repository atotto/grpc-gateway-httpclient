package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/atotto/grpc-gateway-httpclient/test/testdata/apis"
)

var _ apis.MessageServiceServer = (*MessageService)(nil)

type MessageService struct {
	mu       sync.Mutex
	id       int
	messages map[string]*apis.Message
}

func NewMessageService() *MessageService {
	return &MessageService{messages: make(map[string]*apis.Message)}
}

func (s *MessageService) CreateMessage(ctx context.Context, req *apis.CreateMessageRequest) (res *apis.CreateMessageResponse, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id++
	messageID := fmt.Sprintf("%d", s.id)
	s.messages[messageID] = req.GetMessage()
	return &apis.CreateMessageResponse{MessageId: messageID}, nil
}

func (s *MessageService) GetMessage(ctx context.Context, req *apis.GetMessageRequest) (res *apis.GetMessageResponse, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &apis.GetMessageResponse{Message: s.messages[req.GetMessageId()]}, nil
}
