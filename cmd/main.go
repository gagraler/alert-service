package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gagraler/alert-service/internal/controller"
	"github.com/gagraler/alert-service/pkg/cfg"
	"github.com/gagraler/alert-service/pkg/database"
	"github.com/gagraler/alert-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/1/11 21:42
 * @file: main.go
 * @description: 主程序入口
 */

type AlertServiceConfig struct {
	Database database.Config
}

func init() {

	config := &AlertServiceConfig{}

	_, err := cfg.InitCfg("conf.d", "alert-service", "toml", config)
	if err != nil {
		os.Exit(0)
	}

	_, err = database.NewDatabase(&config.Database)
	if err != nil {
		os.Exit(0)
	}
}

func main() {

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()

	controller.InitializeController(g)

	log := logger.SugaredLogger()

	server := http.Server{
		Addr:    "0.0.0.0:8988",
		Handler: g,
	}
	log.Infof("listen: %s", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("server listen err:%s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error: ", err)
	}
	log.Info("Stopped...")
}
