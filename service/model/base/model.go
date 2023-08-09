package base

import (
	"fmt"
	"github.com/capell/capell_scan/lib/app/iapp"
	"github.com/capell/capell_scan/lib/wmongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model struct {
	db       *mongo.Client
	database string
	BaseCol  *wmongo.ColOperator
}

func NewModel(app iapp.IApp, opt *Option) (mdl *Model, err error) {
	mdl = &Model{
		db:       app.GetMongo(opt.Mongo),
		database: opt.DbName,
	}
	if mdl.db == nil {
		err = fmt.Errorf("must provider mongo for base model")
		return
	}
	if len(opt.DbName) == 0 {
		err = fmt.Errorf("must provider mongo database name base model")
		return
	}
	mdl.BaseCol = wmongo.NewColOperator(mdl.db, opt.DbName, "base")
	return
}
