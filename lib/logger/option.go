package logger

type StdoutOption struct {
	Level string `toml:"level"`
}

type FileOption struct {
	Level    string `toml:"level"`
	FileName string `toml:"filename"`
	Format   string `toml:"format"`
	MaxLines int    `toml:"maxlines"`
	MaxSize  int    `toml:"maxsize"`
	NoDaily  bool   `toml:"nodaily"`
	NoRotate bool   `toml:"norotate"`
}

type ScribeOption struct {
	Level         string `toml:"level"`
	Endpoint      string `toml:"endpoint"`
	Category      string `toml:"category"`
	Format        string `toml:"format"`
	SuffixEnabled bool   `toml:"suffix_enabled"`
}

type DingDingOption struct {
	Level     string   `toml:"level"`
	Key       string   `toml:"key"`
	Url       string   `toml:"url"`
	AtMobiles []string `toml:"at_mobiles"`
}

type Option struct {
	Level          string                     `toml:"level"`
	Stdout         *StdoutOption              `toml:"stdout"`
	Files          map[string]*FileOption     `toml:"file"`
	Scribes        map[string]*ScribeOption   `toml:"scribe"`
	Dingdings      map[string]*DingDingOption `toml:"dingding"`
	DisableMetrics bool                       `toml:"disable_metrics"`
}
