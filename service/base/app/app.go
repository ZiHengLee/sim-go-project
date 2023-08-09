package app

import (
	"fmt"
	"github.com/capell/capell_scan/lib/app"
	"github.com/capell/capell_scan/lib/discovery"
	"github.com/capell/capell_scan/lib/logger"
	"github.com/capell/capell_scan/service/base/handler"
	"github.com/capell/capell_scan/service/worker"
	"log"
)

type App struct {
	app.App

	opt Option
}

func (a App) Start() {
	err := a.App.Init(&a.opt.Option, &a.opt)
	if err != nil {
		log.Fatalf("init app err:%v", err)
		return
	}
	a.opt.Handler.App = &a
	h, err := handler.NewHandler(a.opt.Handler)
	if err != nil {
		logger.Error("init handler err:%v", err)
		return
	}
	path := "/" + a.opt.Handler.Name
	h.Route(a.HttpServer().Engine().Group(path))

	if a.opt.Etcd != nil {
		etcdUrl := a.opt.Etcd.Addr
		etcdRegister := discovery.NewRegister([]string{etcdUrl})
		defer etcdRegister.Stop()
		baseNode := discovery.Server{
			Name: a.opt.Handler.Name,
			Addr: a.opt.GrpcServer.Addr,
		}
		if _, err := etcdRegister.Register(baseNode, 10); err != nil {
			panic(fmt.Sprintf("start server failed, err: %v", err))
		}
		logger.Info("init etcd success")
	}

	err = worker.Init(&a, &a.opt.Worker)
	if err != nil {
		logger.Error("init worker err:%v", err)
		return
	}
	logger.Info("init worker success")
	go worker.Run()
	a.Run()
}
