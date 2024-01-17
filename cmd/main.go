package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/keington/alertService/internel/controller"
	"github.com/keington/alertService/pkg/cfg"
	"github.com/keington/alertService/pkg/database"
	"os"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/11 21:42
 * @file: main.go
 * @description: 主程序入口
 */

type AlertServiceConfig struct {
	Database database.Config
}

func (c *AlertServiceConfig) LoadConfigStruct() {
}

func init() {
	flag.Parse()

	config := &AlertServiceConfig{}

	_, err := cfg.InitCfg("./conf.d", "alertService", "toml", config)
	if err != nil {
		os.Exit(0)
	}

	_, err = database.NewDatabase(&config.Database)
	if err != nil {
		os.Exit(0)
	}

	return
}

func main() {

	g := gin.Default()

	controller.InitializeController(g)

	err := g.Run("0.0.0.0:8588")
	if err != nil {
		os.Exit(0)
	}
}
