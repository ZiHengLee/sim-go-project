package app

import (
	"github.com/capell/capell_scan/lib/app"
	"github.com/capell/capell_scan/service/base/handler"
	"github.com/capell/capell_scan/service/worker"
)

type Option struct {
	app.Option

	//Model   model.Option   `toml:"model"`
	Worker  worker.Option   `toml:"worker"`
	Handler *handler.Option `toml:"handler"`
}
