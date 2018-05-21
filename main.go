// Package main v1接口.
//
// 接口文档
//
//     Schemes: http
//     Host: cxz.125i.cn
//     BasePath: /v1
//     Version: 0.0.1
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/hfdend/cxz/handler/v1"
)

func main() {
	cli.Init()
	engine := gin.Default()
	route(engine)
	log.Printf("server run %s\n", conf.Config.Main.Addr)
	log.Fatalln(engine.Run(conf.Config.Main.Addr))
}

func route(engine *gin.Engine) {
	{
		var MustLogin = v1.MustLogin
		g := engine.Group("v1", v1.SetLoginUser)

		g.POST("register/send", v1.Passport.RegisterSend)
		g.POST("register", v1.Passport.Register)
		g.POST("login", v1.Passport.Login)
		g.GET("user", v1.User.GetUserInfo)

		g.GET("district/gradation", v1.District.GetGradation)
		g.POST("address/save", MustLogin, v1.Address.Save)
		g.POST("address/del", MustLogin, v1.Address.Del)
		g.POST("address/list", MustLogin, v1.Address.List)
	}
}
