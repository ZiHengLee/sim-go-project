package handler

import (
	"github.com/capell/capell_scan/gateway/handler/base"
	"github.com/capell/capell_scan/lib/app/iapp"
)

type Option struct {
	Base *base.Option `toml:"base"`

	App iapp.IApp `toml:"-"`
}
