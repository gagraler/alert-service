package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keington/alertService/internel/controller"
	"github.com/keington/alertService/pkg/cfg"
	"github.com/keington/alertService/pkg/database"
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

	server := http.Server{
		Addr:    "0.0.0.0:8588",
		Handler: g,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server listen err:%s", err)
		}
	}()

	//err := g.Run("0.0.0.0:8588")
	//if err != nil {
	//	os.Exit(0)
	//}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error: ", err)
	}
	log.Println("service closing...")
}
