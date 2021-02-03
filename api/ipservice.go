package main

import (
	"flag"
	"fmt"

	"ip-service/api/internal/config"
	"ip-service/api/internal/handler"
	"ip-service/api/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "etc/ipService.yaml", "the config file")
var dbFile = flag.String("d", "https://data-oss-hewei.oss-cn-shenzhen.aliyuncs.com/database/ip2region.db", "the ip database file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if dbFile != nil && (*dbFile) != "" {
		c.Database = *dbFile
	}
	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	handler.RegisterHandlers(server, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
