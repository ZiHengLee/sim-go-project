package worker

import (
	"github.com/capell/capell_scan/app/service/worker/hello_world"
)

type Option struct {
	HelloWorld *hello_world.Option `toml:"hello_world"`
}
