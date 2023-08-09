package model

import (
	"github.com/capell/capell_scan/service/model/base"
)

type Option struct {
	Base *base.Option `toml:"base"`
}
