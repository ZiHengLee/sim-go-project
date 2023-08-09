package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/capell/capell_scan/lib/logger"
	"time"
)

func callPanic() {
	defer func() {
		if err := recover(); err != nil {
			logger.Alert("panic with err:%v", err)
		}
	}()

	panic("direct p a n i c")
}

func main() {
	var err error
	var opt logger.Option
	fmt.Printf("logger example\n")
	_, err = toml.DecodeFile("example.toml", &opt)
	if err != nil {
		fmt.Printf("load option file err:%v\n", err)
		return
	}
	err = logger.Init(&opt)
	if err != nil {
		fmt.Printf("init logger err:%v\n", err)
		return
	}
	callPanic()
	for i := 0; i < 10; i++ {
		fmt.Printf("iter:%v\n", i)
		logger.Debug("debug log i=%v", i)
		logger.Trace("trace log i=%v", i)
		logger.Info("info log i=%v", i)
		logger.Warn("warn log i=%v", i)
		logger.Error("error log=%v", i)
	}

	time.Sleep(time.Second)
	logger.Close()
}
