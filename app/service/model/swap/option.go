package swap

type Option struct {
	Mongo  string `toml:"mongo"`
	DbName string `toml:"db_name"`
	Sql    string `toml:"sql"`
}
