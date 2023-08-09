package handler

import (
	"github.com/capell/capell_scan/lib/app/iapp"
)

type Option struct {
	App      iapp.IApp `toml:"-"`
	Name     string    `toml:"name"`
	Etcd     bool      `toml:"etcd"`
	NeedGrpc bool      `toml:"need_grpc"`
}
