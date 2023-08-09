package main

import (
	"github.com/capell/capell_scan/service/base/app"
	"time"
)

func main() {
	app.App{}.Start()
	time.Sleep(time.Second) //wait logger
}
