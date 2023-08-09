package model

import (
	"github.com/capell/capell_scan/lib/app/iapp"
	"github.com/capell/capell_scan/lib/logger"
	"github.com/capell/capell_scan/service/model/base"
)

var (
	BaseMdl *base.Model
)

func Init(app iapp.IApp, opt *Option) (err error) {
	if opt.Base != nil {
		m, err := base.NewModel(app, opt.Base)
		logger.Info("new base model opt:%#v err:%v", opt.Base, err)
		if err != nil {
			return err
		}
		BaseMdl = m
	}
	return
}
