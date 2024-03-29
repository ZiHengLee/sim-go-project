package app

import (
	"github.com/capell/capell_scan/gateway/handler"
	"github.com/capell/capell_scan/lib/app"
)

type Option struct {
	app.Option

	//Model   model.Option   `toml:"model"`
	//Worker  worker.Option  `toml:"worker"`
	Handler *handler.Option `toml:"handler"`
	//
	//EthClient *ethclient.Option `toml:"eth_client"`
}
