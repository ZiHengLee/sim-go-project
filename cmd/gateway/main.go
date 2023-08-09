package main

import (
	"github.com/capell/capell_scan/app/gateway/app"
	"time"
)

func main() {
	app.App{}.Start()
	time.Sleep(time.Second) //wait logger
}
