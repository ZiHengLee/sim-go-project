package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

type ResType struct {
	Detail  string `json:"detail"`
	JsonRPC string `json:"jsonrpc"`
	Id      int64
	Error   interface{} `json:"error"`
}

type CustomReq struct {
}

func (c *CustomReq) NewRequest(ctx context.Context, host, path string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", host+"/"+path, nil)
	if req != nil {
		req.Header.Add("X-Forward-For", "192.168.1.1")
	}
	fmt.Printf("create custom request:%#v\n", req)
	return req, err
}

func TestJsonCall(t *testing.T) {
	opt := Option{
		Host: "https://rpc.theindex.io",
	}
	cli, err := NewHttpClient(&opt)
	if err != nil {
		t.Errorf("new httpclient err:%v", err)
	}
	res, err := JsonRpc[ResType](cli, context.Background(), "", GetReq)
	if err != nil {
		t.Errorf("json rpc with GET method err:%v", err)
	}
	fmt.Printf("GET result:%#v\n", res)

	custom := &CustomReq{}
	res, err = JsonRpc[ResType](cli, context.Background(), "", custom)
	if err != nil {
		t.Errorf("json rpc with GET method err:%v", err)
	}
	fmt.Printf("Custom request result:%#v\n", res)

	param := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
	}
	res, err = JsonRpc[ResType](cli, context.Background(), "", param)
	if err != nil {
		t.Errorf("jsoncall err:%v", err)
	}
	fmt.Printf("jsoncall res:%#v\n", res)

	res, err = AsyncJsonRpc[ResType](cli, context.Background(), "", param).Get()
	if err != nil {
		t.Errorf("async jsoncall err:%v", err)
	}
	fmt.Printf("async jsoncall res:%#v\n", res)

}
