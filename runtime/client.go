package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/schema"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Client struct {
	Client       *http.Client
	DumpRequest  func(req *http.Request)
	DumpResponse func(resp *http.Response)
}

var DefaultQueryEncoder = schema.NewEncoder()

func init() {
	DefaultQueryEncoder.SetAliasTag("json")
}

var DefaultClient = &Client{
	Client: http.DefaultClient,
}

func Get(ctx context.Context, u string, request, response interface{}) error {
	return DefaultClient.Get(ctx, u, request, response)
}

func Post(ctx context.Context, u string, request, response interface{}) error {
	return DefaultClient.Post(ctx, u, request, response)
}

func Patch(ctx context.Context, u string, request, response interface{}) error {
	return DefaultClient.Patch(ctx, u, request, response)
}

func Delete(ctx context.Context, u string, request, response interface{}) error {
	return DefaultClient.Delete(ctx, u, request, response)
}

// ---

func (c *Client) Get(ctx context.Context, u string, request, response interface{}) error {
	u = toQueryURL(DefaultQueryEncoder, u, request)
	req, err := NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	return c.Do(ctx, req, response)
}

func (c *Client) Post(ctx context.Context, u string, request, response interface{}) error {
	req, err := NewRequest("POST", u, request)
	if err != nil {
		return err
	}
	return c.Do(ctx, req, response)
}

func (c *Client) Patch(ctx context.Context, u string, request, response interface{}) error {
	req, err := NewRequest("PATCH", u, request)
	if err != nil {
		return err
	}
	return c.Do(ctx, req, response)
}

func (c *Client) Delete(ctx context.Context, u string, request, response interface{}) error {
	u = toQueryURL(DefaultQueryEncoder, u, request)
	req, err := NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	return c.Do(ctx, req, response)
}

// ---

func NewRequest(method, url string, request interface{}) (*http.Request, error) {
	if request == nil || reflect.ValueOf(request).IsNil() {
		return http.NewRequest(method, url, nil)
	}

	body := &bytes.Buffer{}
	switch v := request.(type) {
	case proto.Message:
		marshaler := jsonpb.Marshaler{}
		if err := marshaler.Marshal(body, v); err != nil {
			return nil, fmt.Errorf("json encode: %s", err)
		}
	default:
		if err := json.NewEncoder(body).Encode(request); err != nil {
			return nil, fmt.Errorf("json encode: %s", err)
		}
	}

	return http.NewRequest(method, url, body)
}

type Status struct {
	spb.Status
	ErrorMessage string `json:"error,omitempty"`
}

func (c *Client) Do(ctx context.Context, req *http.Request, response interface{}) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		for key, values := range md {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	req.Header.Set("Content-Type", "application/json")

	if c.DumpRequest != nil {
		c.DumpRequest(req)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c.DumpResponse != nil {
		c.DumpResponse(resp)
	}

	if resp.StatusCode >= 400 {
		st := &Status{}
		err := json.NewDecoder(resp.Body).Decode(st)
		if err != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("json decode error status=%s, err=%s, body=%s", resp.Status, err, string(body))
		}
		if st.Message == "" {
			st.Message = st.ErrorMessage
		}
		return status.ErrorProto(&st.Status)
	}

	dec := json.NewDecoder(resp.Body)
	switch v := response.(type) {
	case proto.Message:
		err = (&jsonpb.Unmarshaler{AllowUnknownFields: true}).UnmarshalNext(dec, v)
	default:
		err = dec.Decode(response)
	}
	if err != nil {
		return fmt.Errorf("type %T http %s", response, err)
	}
	return nil
}
