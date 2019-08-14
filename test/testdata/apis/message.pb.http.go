// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: message.proto

package apis

import (
	context "context"
	fmt "fmt"
	runtime "github.com/atotto/grpc-gateway-httpclient/runtime"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	math "math"
	http "net/http"
	url "net/url"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func NewMessageServiceHttpClient(client *http.Client, url *url.URL) *MessageServiceHttpClient {
	return &MessageServiceHttpClient{Client: &runtime.Client{Client: client}, URL: *url}
}

type MessageServiceHttpClient struct {
	Client *runtime.Client
	URL    url.URL
}

func (c *MessageServiceHttpClient) CreateMessage(ctx context.Context, req *CreateMessageRequest, opts ...grpc.CallOption) (*CreateMessageResponse, error) {
	var err error
	res := &CreateMessageResponse{}
	u := c.URL
	u.Path = "/messages"
	err = c.Client.Post(ctx, u.String(), req, res)
	return res, err
}

func (c *MessageServiceHttpClient) GetMessage(ctx context.Context, req *GetMessageRequest, opts ...grpc.CallOption) (*GetMessageResponse, error) {
	var err error
	res := &GetMessageResponse{}
	u := c.URL
	u.Path = "/messages/" + req.MessageId + ""
	err = c.Client.Get(ctx, u.String(), req, res)
	return res, err
}