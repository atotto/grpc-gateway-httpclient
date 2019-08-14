package generate

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/atotto/grpc-gateway-httpclient/test/service"
	"github.com/atotto/grpc-gateway-httpclient/test/testdata/apis"

	"github.com/akutz/memconn"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func TestHttpClient(t *testing.T) {
	grpcServer := grpc.NewServer()

	apis.RegisterEchoServiceServer(grpcServer, service.NewEchoService())
	apis.RegisterMessageServiceServer(grpcServer, service.NewMessageService())

	listener, err := memconn.Listen("memu", "grpc")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			if e, ok := err.(*net.OpError); ok {
				if e.Err.Error() == "listener closed" {
					return
				}
			}
			t.Fatal(err)
		}
	}()

	conn, err := grpc.Dial(
		"grpc",
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return memconn.DialContext(ctx, "memu", addr)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	mux := runtime.NewServeMux()
	apis.RegisterEchoServiceHandler(ctx, mux, conn)
	apis.RegisterMessageServiceHandler(ctx, mux, conn)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	//----

	u, _ := url.Parse(ts.URL)

	t.Run("echo service", func(t *testing.T) {
		client := apis.NewEchoServiceHttpClient(http.DefaultClient, u)

		client.Client.DumpRequest = func(req *http.Request) {
			dump, _ := httputil.DumpRequestOut(req, true)
			t.Logf("%q", dump)
		}
		client.Client.DumpResponse = func(resp *http.Response) {
			dump, _ := httputil.DumpResponse(resp, true)
			t.Logf("%q", dump)
		}

		res, err := client.Echo(ctx, &apis.EchoRequest{Echo: &apis.Echo{Message: "hello"}})
		if err != nil {
			t.Fatal(err)
		}

		if msg := res.GetEcho().GetMessage(); msg != "hello" {
			t.Fatalf("want hello, got %s", msg)
		}
	})

	t.Run("message service", func(t *testing.T) {
		client := apis.NewMessageServiceHttpClient(http.DefaultClient, u)

		client.Client.DumpRequest = func(req *http.Request) {
			dump, _ := httputil.DumpRequestOut(req, true)
			t.Logf("%q", dump)
		}
		client.Client.DumpResponse = func(resp *http.Response) {
			dump, _ := httputil.DumpResponse(resp, true)
			t.Logf("%q", dump)
		}

		{
			res, err := client.CreateMessage(ctx, &apis.CreateMessageRequest{
				Message: &apis.Message{Content: "hello"},
			})
			if err != nil {
				t.Fatal(err)
			}

			if msgID := res.GetMessageId(); msgID != "1" {
				t.Fatalf("want 1, got %s", msgID)
			}
		}
		{
			res, err := client.GetMessage(ctx, &apis.GetMessageRequest{
				MessageId: "1",
			})
			if err != nil {
				t.Fatal(err)
			}

			if content := res.GetMessage().GetContent(); content != "hello" {
				t.Fatalf("want 1, got %s", content)
			}
		}
	})

}
