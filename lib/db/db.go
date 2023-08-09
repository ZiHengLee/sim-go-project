package db

import (
	"fmt"

	"github.com/capell/capell_scan/lib/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbs map[string]*gorm.DB

func Init(opts map[string]Option) (err error) {
	if dbs != nil {
		err = fmt.Errorf("db package already init")
		return
	}
	dbs = make(map[string]*gorm.DB, len(opts))
	for k, opt := range opts {
		db, err := newDB(opt)
		if err != nil {
			return err
		}
		dbs[k] = db
	}

	if err != nil {
		logger.Error("register metric:db_latency err:%v", err)
		return err
	}

	return
}

func GetDB(key string) (db *gorm.DB, err error) {
	if dbs == nil {
		err = fmt.Errorf("db package is not init")
		return
	}
	db = dbs[key]
	if db == nil {
		err = fmt.Errorf("db:%v is not exists", key)
	}
	return
}

func newDB(opt Option) (db *gorm.DB, err error) {
	var d gorm.Dialector
	if opt.Driver == "sqlite" {
		d = sqlite.Open(opt.Url)
	} else {
		d = mysql.Open(opt.Url)
	}
	logger.Info("open database:%#v", opt)
	lg := NewGormLogger(&opt)
	cfg := &gorm.Config{
		Logger: lg,
	}
	db, err = gorm.Open(d, cfg)
	if err != nil {
		logger.Error("open database:%#v err:%v", opt, err)
		return nil, err
	}
	return
}
