package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxr/cli"
	"github.com/hfdend/cxr/conf"
)

func main() {
	cli.Init()
	engine := gin.Default()
	log.Printf("server run %s%\n", conf.Config.Main.Addr)
	log.Fatalln(engine.Run(conf.Config.Main.Addr))
}
