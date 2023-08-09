package worker

import (
	"github.com/capell/capell_scan/lib/app/iapp"
	"github.com/capell/capell_scan/lib/logger"
	"github.com/capell/capell_scan/service/worker/hello_world"
)

var (
	HelloWorker *hello_world.Worker
)

func Init(app iapp.IApp, opt *Option) (err error) {
	if opt.HelloWorld != nil {
		w, err := hello_world.NewWorker(app, opt.HelloWorld)
		logger.Info("new common worker opt:%#v err:%v", opt.HelloWorld, err)
		if err != nil {
			logger.Error("new worker for common err:%v", err)
			return err
		}
		HelloWorker = w
		hello_world.DefaultWorker = w
	}
	return
}

func Run() {
	if HelloWorker != nil {
		go HelloWorker.Run()
	}
}
