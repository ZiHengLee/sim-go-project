package httpclient

type Option struct {
	Host      string `toml:"host"`
	Proxy     string `toml:"proxy"`
	NeedProxy bool   `toml:"need_proxy"`
}

type PkgOption struct {
	Metrics string             `json:"metrics"`
	Clients map[string]*Option `json:"clients"`
}
