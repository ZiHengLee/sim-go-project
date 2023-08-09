package httpclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type replyBody struct {
	err error
	res *http.Response
}

type Reply struct {
	body *replyBody
	ch   chan *replyBody
}

func NewReply(err error) (reply *Reply) {
	reply = &Reply{
		body: &replyBody{
			err: err,
		},
	}
	return
}

func (r *Reply) wait() {
	if r.ch == nil {
		return
	}
	select {
	case b, ok := <-r.ch:
		if ok {
			r.body = b
		}
	}
}

func (r *Reply) Err() (err error) {
	r.wait()
	err = r.body.err
	return
}

func (r *Reply) Get() (res *http.Response, err error) {
	r.wait()
	res, err = r.body.res, r.body.err
	return
}

func (r *Reply) StatusCode() (code int) {
	r.wait()
	code = r.body.res.StatusCode
	return
}

func (r *Reply) GetJson(v interface{}) (res *http.Response, err error) {
	res, err = r.Get()
	if err != nil {
		return
	}
	body := res.Body
	defer body.Close()
	dat, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(dat, v)
	return
}
