package model

import (
	"github.com/capell/capell_scan/app/service/model/swap"
)

type Option struct {
	Swap *swap.Option `toml:"swap"`
}
