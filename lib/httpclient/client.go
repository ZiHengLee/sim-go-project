package httpclient

import (
	"net/http"
	"net/url"
)

type HttpClient struct {
	opt    *Option
	client *http.Client
}

func NewHttpClient(opt *Option) (c *HttpClient, err error) {
	cli := &http.Client{}
	c = &HttpClient{
		opt:    opt,
		client: cli,
	}
	return
}

func NewProxyHttpClient(opt *Option) (c *HttpClient, err error) {
	cli := new(http.Client)
	uri, _ := url.Parse(opt.Proxy)
	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyURL(uri),
	}
	cli.Transport = transport
	c = &HttpClient{
		opt:    opt,
		client: cli,
	}
	return
}

func GetHttpClient(name string) (c *HttpClient) {
	c = clients[name]
	return
}

func (c HttpClient) Send(req *http.Request) (reply *Reply) {
	reply = &Reply{
		ch: make(chan *replyBody, 1),
	}
	go func() {
		//var status string
		//b := time.Now()

		res, err := c.client.Do(req)
		body := &replyBody{
			err: err,
			res: res,
		}

		//if err != nil {
		//	status = "error"
		//} else if res != nil {
		//	status = fmt.Sprintf("%v", res.StatusCode)
		//}
		tgt := c.opt.Host
		if len(tgt) == 0 {
			tgt = req.Host
		}
		if len(tgt) == 0 {
			tgt = req.URL.Host
		}

		//if reqMetric != nil {
		//	reqMetric.Add(1, tgt, req.URL.Path, status)
		//}
		//if latencyMetric != nil {
		//	elapsed := time.Now().Sub(b)
		//	latencyMetric.Observe(float64(elapsed)/float64(time.Millisecond), tgt, req.URL.Path)
		//}

		reply.ch <- body
		close(reply.ch)
	}()
	return
}
