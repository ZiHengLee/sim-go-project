package base

import "github.com/capell/capell_scan/lib/app/iapp"

type Option struct {
	App iapp.IApp `toml:"-"`
}
