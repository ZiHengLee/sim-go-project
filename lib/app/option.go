package app

import (
	"github.com/capell/capell_scan/lib/db"
	"github.com/capell/capell_scan/lib/httpclient"
	"github.com/capell/capell_scan/lib/httpserver"
	"github.com/capell/capell_scan/lib/logger"
)

type TLSOption struct {
	MinVersion         uint16 `toml:"min_version"`
	MaxVersion         uint16 `toml:"max_version"`
	ServerName         string `toml:"server_name"`
	InsecureSkipVerify bool   `toml:"insecure_skip_verify"`
	CAFile             string `toml:"ca_file"`
}

type DbOption struct {
	Driver string `toml:"driver"`
	Url    string `toml:"url"`
}

type RedisOption struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
	TLS      string `toml:"tls"`
}

type MongoOption struct {
	Url string `toml:"url"`
	TLS string `toml:"tls"`
}

type Option struct {
	Log        logger.Option          `toml:"log"`
	HttpServer *httpserver.Option     `toml:"httpserver"`
	GrpcServer *GrpcOption            `toml:"grpcserver"`
	HttpClient *httpclient.PkgOption  `toml:"httpclient"`
	Databases  map[string]db.Option   `toml:"databases"`
	Redis      map[string]RedisOption `toml:"redis"`
	Mongos     map[string]MongoOption `toml:"mongos"`
	Etcd       *EtcdOption            `toml:"etcd"`
}

type GrpcOption struct {
	Addr  string   `toml:"addr"`
	Nodes []string `toml:"nodes"`
}

type EtcdOption struct {
	Addr string `toml:"addr"`
}
