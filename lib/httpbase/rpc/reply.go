package rpc

import (
	"encoding/json"
	"io/ioutil"

	"github.com/capell/capell_scan/lib/httpclient"
)

type Reply struct {
	reply *httpclient.Reply
}

func (r Reply) Get(resp interface{}) (err error) {
	res, err := r.reply.Get()
	if err != nil {
		return
	}

	body := res.Body
	defer body.Close()
	dat, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(dat, resp)
	return
}
