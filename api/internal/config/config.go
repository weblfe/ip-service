package config

import "github.com/tal-tech/go-zero/rest"

type Config struct {
	rest.RestConf
	Database string `json:"Database,optional,default=db/ip2region.db"`
	TempDir  string `json:"TempDir,optional,default=tmp/"`
	DocsPath  string `json:",optional,default=./docs/swagger.json"`
}
