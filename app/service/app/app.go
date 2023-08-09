package app

import (
	"fmt"
	"github.com/capell/capell_scan/app/service/handler"
	"github.com/capell/capell_scan/app/service/worker"
	"github.com/capell/capell_scan/lib/app"
	"github.com/capell/capell_scan/lib/discovery"
	"github.com/capell/capell_scan/lib/logger"
	"github.com/capell/capell_scan/rpc"
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
	err = handler.Init(a.opt.Handler)
	if err != nil {
		logger.Error("init handler err:%v", err)
		return
	}
	rpc.Init()
	etcdRegister := discovery.NewRegister([]string{"127.0.0.1:2379"})
	defer etcdRegister.Stop()
	taskNode := discovery.Server{
		Name: "swap",
		Addr: "127.0.0.1:10002",
	}
	if _, err := etcdRegister.Register(taskNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logger.Info("init worker")
	err = worker.Init(&a, &a.opt.Worker)
	if err != nil {
		logger.Error("init worker err:%v", err)
		return
	}
	go worker.Run()
	handler.Route(a.HttpServer().Engine())
	a.Run()
}
