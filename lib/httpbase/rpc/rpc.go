package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/capell/capell_scan/lib/httpclient"
	"net/http"
	"strings"
)

type Rpc struct {
	opt    *Option
	client *httpclient.HttpClient
}

func NewRpc(opt *Option) (c *Rpc, err error) {
	cli, err := httpclient.NewHttpClient(&opt.Option)
	if err != nil {
		return nil, err
	}
	c = &Rpc{
		opt: opt, client: cli,
	}
	return
}

func (r *Rpc) Send(ctx context.Context, path string, param interface{}) (reply *Reply) {
	reply = &Reply{}
	sep := "/"
	if strings.HasPrefix(path, "/") {
		sep = ""
	}
	url := fmt.Sprintf("%s%s%s", r.opt.Host, sep, path)
	dat, err := json.Marshal(param)
	if err != nil {
		reply.reply = httpclient.NewReply(err)
		return
	}
	buf := bytes.NewReader(dat)
	req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		reply.reply = httpclient.NewReply(err)
		return
	}
	reply.reply = r.client.Send(req)
	return
}
