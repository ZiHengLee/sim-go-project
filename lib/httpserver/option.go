package httpserver

type Option struct {
	Addr              string `toml:"addr"`
	EnableDebug       bool   `toml:"enable_debug"`
	EnableHealthCheck bool   `toml:"enable_healthcheck"`
	DisableMetrics    bool   `toml:"disable_metrics"`
	HystrixName       string `toml:"hystrix_name"`
}
