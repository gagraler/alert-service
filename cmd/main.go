package main

import (
	"github.com/gin-gonic/gin"
	"github.com/keington/alart-service/internel/controller"
	"os"
)

/**
 * @author: x.gallagher.anderson@gmail.com
 * @time: 2024/1/11 21:42
 * @file: main.go
 * @description: 主程序入口
 */

func init() {

}

func main() {

	g := gin.Default()

	controller.InitializeController(g)

	err := g.Run(":8588")
	if err != nil {
		os.Exit(0)
	}
}
