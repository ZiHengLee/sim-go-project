package httpclient

import (
	"github.com/capell/capell_scan/lib/logger"
)

var clients map[string]*HttpClient

func Init(opt *PkgOption) (err error) {
	if len(opt.Clients) > 0 {
		clis := make(map[string]*HttpClient)
		for name, copt := range opt.Clients {
			var c *HttpClient
			if copt.NeedProxy {
				c, err = NewProxyHttpClient(copt)
				if err != nil {
					logger.Error("new httpclient with opt:%#v err:%v", copt, err)
					return err
				}
			} else {
				c, err = NewHttpClient(copt)
				if err != nil {
					logger.Error("new httpclient with opt:%#v err:%v", copt, err)
					return err
				}
			}
			clis[name] = c
		}
		clients = clis
	}
	return
}
