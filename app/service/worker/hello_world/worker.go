package hello_world

import (
	"github.com/capell/capell_scan/lib/app/iapp"
	"github.com/capell/capell_scan/lib/logger"
)

type Worker struct {
	opt *Option
}

var DefaultWorker *Worker

func NewWorker(app iapp.IApp, opt *Option) (w *Worker, err error) {
	w = &Worker{
		opt: opt,
	}
	return
}

func (w *Worker) Run() {
	logger.Info("hello")
}
