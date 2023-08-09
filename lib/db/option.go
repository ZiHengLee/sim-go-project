package db

type Option struct {
	Name   string `toml:"name"`
	Driver string `toml:"driver"`
	Url    string `toml:"url"`

	SlowThreshold float64 `toml:"slow_threshold"`
}
