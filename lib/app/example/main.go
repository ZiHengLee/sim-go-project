package main

import (
	"log"
	"time"

	"github.com/capell/capell_scan/app"
	"github.com/gin-gonic/gin"
)

type Option struct {
	app.Option
}

type App struct {
	app.App

	opt Option
}

func main() {
	a := App{}
	err := a.Init(&a.opt.Option, &a.opt)
	if err != nil {
		log.Fatalf("init err:%v", err)
		return
	}
	e := a.HttpServer().Engine()
	e.GET("/now", func(c *gin.Context) {
		c.String(200, "now:%v", time.Now())
	})
	a.Run()
}
