package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type INewHttpRequest interface {
	NewRequest(ctx context.Context, host, path string, body any) (*http.Request, error)
}

type IMethod interface {
	RpcMethod() string
}

func RpcMethod(param any) string {
	if im, ok := param.(IMethod); ok {
		return im.RpcMethod()
	}
	return "POST"
}

type RpcGetReq struct {
}

func (RpcGetReq) RpcMethod() string {
	return "GET"
}

var GetReq RpcGetReq

func rpc[ResType any, ReqType any](cli *HttpClient, ctx context.Context, path string, param ReqType) (reply *Reply, err error) {
	var req *http.Request
	var iparam interface{}
	iparam = param
	if inr, ok := iparam.(INewHttpRequest); ok {
		req, err = inr.NewRequest(ctx, cli.opt.Host, path, param)
	} else {
		url := cli.opt.Host + "/" + path
		var body io.Reader
		method := RpcMethod(param)
		if method != "GET" {
			dat, err := json.Marshal(param)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(dat)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, body)
	}
	if err != nil {
		return nil, err
	}
	reply = cli.Send(req)
	return
}

func JsonRpc[ResType any, ReqType any](cli *HttpClient, ctx context.Context, path string, param ReqType) (res *ResType, err error) {
	reply, err := rpc[ResType, ReqType](cli, ctx, path, param)
	if err != nil {
		return nil, err
	}
	var r ResType
	_, err = reply.GetJson(&r)
	if err == nil {
		res = &r
	}
	return
}

type JsonRpcReply[ResType any] struct {
	reply *Reply
}

func (r JsonRpcReply[ResType]) Get() (res *ResType, err error) {
	var ret ResType
	_, err = r.reply.GetJson(&ret)
	if err == nil {
		res = &ret
	}
	return
}

func (r JsonRpcReply[ResType]) Err() (err error) {
	return r.reply.Err()
}

func (r JsonRpcReply[ResType]) HttpCode() (httpCode int) {
	return r.reply.StatusCode()
}

func AsyncJsonRpc[ResType any, ReqType any](cli *HttpClient, ctx context.Context, path string, param ReqType) JsonRpcReply[ResType] {
	reply, err := rpc[ResType, ReqType](cli, ctx, path, param)
	if err != nil {
		reply = NewReply(err)
	}
	return JsonRpcReply[ResType]{reply}
}
