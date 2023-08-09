package app

import (
	"github.com/capell/capell_scan/gateway/handler"
	"github.com/capell/capell_scan/lib/app"
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
	err = handler.Init(a.opt.Handler)
	if err != nil {
		logger.Error("init handler err:%v", err)
		return
	}
	rpc.Init()
	handler.Route(a.HttpServer().Engine())
	a.Run()
}
