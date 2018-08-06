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
	"github.com/hfdend/cxz/handler"
	"github.com/hfdend/cxz/handler/api"
	"github.com/hfdend/cxz/handler/v1"
	"github.com/hfdend/cxz/role"
)

func main() {
	cli.Init()
	engine := gin.Default()
	route(engine)
	log.Printf("server run %s\n", conf.Config.Main.Addr)
	log.Fatalln(engine.Run(conf.Config.Main.Addr))
}

func route(engine *gin.Engine) {
	engine.POST("wxpayment/notify", handler.WXAPaymentNotify)
	{
		var MustLogin = v1.MustLogin
		g := engine.Group("v1", v1.SetLoginUser)
		g.POST("miniprogram/login", v1.Passport.LoginByJsCode)
		g.POST("bind/phone", v1.Passport.BindPhone)
		g.POST("register/send", v1.Passport.RegisterSend)
		g.POST("register", v1.Passport.Register)
		g.POST("login", v1.Passport.Login)

		g.GET("district/gradation", v1.District.GetGradation)
		g.GET("user", MustLogin, v1.User.GetUserInfo)
		g.POST("address/save", MustLogin, v1.Address.Save)
		g.POST("address/del", MustLogin, v1.Address.Del)
		g.GET("address/list", MustLogin, v1.Address.List)
		g.GET("address/detail", MustLogin, v1.Address.GetByID)

		g.GET("product/attribute/items", MustLogin, v1.Product.AttributeItems)
		g.GET("product/list", MustLogin, v1.Product.GetList)
		g.GET("product/detail", MustLogin, v1.Product.GetByID)

		g.POST("order/products/freight", MustLogin, v1.Order.GetFreight)
		g.POST("order/build", MustLogin, v1.Order.Build)
		g.GET("order/detail", MustLogin, v1.Order.GetByOrderID)
		g.POST("order/list", MustLogin, v1.Order.GetList)
		g.POST("order/wxapayment", MustLogin, v1.Order.WXAPayment)
		g.GET("order/plans", MustLogin, v1.Order.GetOrderPlanList)
		g.GET("order/query/express", MustLogin, v1.Order.QueryExpress)
		g.POST("order/plan/delay", MustLogin, v1.Order.PlanDelay)

		g.GET("banner/list", v1.Banner.GetList)
	}
	{
		var MustLogin = api.Passport.MustLogin
		g := engine.Group("api", api.Passport.SetLoginUser)
		{
			g.POST("file/upload", MustLogin(), api.File.Upload)
		}
		{
			g.GET("admin_user/group", MustLogin([]string{role.AdminUserGroup}), api.AdminUser.GetGroupList)
			g.GET("admin_user/group/info", MustLogin([]string{role.AdminUserGroupInfo}), api.AdminUser.GetGroupByID)
			g.POST("admin_user/group/save", MustLogin([]string{role.AdminUserGroupEdit}), api.AdminUser.GroupSave)
			g.POST("admin_user/group/del", MustLogin([]string{role.AdminUserGroupDel}), api.AdminUser.GroupDel)
			g.GET("admin_user/users", MustLogin([]string{role.AdminUser}), api.AdminUser.GetUserList)
			g.GET("admin_user/user/info", MustLogin([]string{role.AdminUserEdit}), api.AdminUser.GetUser)
			g.POST("admin_user/user/save", MustLogin([]string{role.AdminUserEdit}), api.AdminUser.SaveUser)
			g.POST("admin_user/user/del", MustLogin([]string{role.AdminUserDel}), api.AdminUser.DelUser)
			g.GET("admin_user/role/list", MustLogin([]string{role.AdminUserRole}), api.AdminUser.RoleList)

			g.POST("passport/update/password", MustLogin(), api.Passport.UpdatePassword)
			g.POST("passport/login", api.Passport.Login)
			g.POST("passport/login/out", MustLogin(), api.Passport.LoginOut)
			g.GET("passport/user", MustLogin(), api.Passport.GetUser)
		}
		{
			g.GET("attribute/list", MustLogin(), api.Attribute.GetList)
			g.GET("attribute/list/detail", MustLogin(), api.Attribute.GetAll)
			g.POST("attribute/items/save", MustLogin(), api.Attribute.SaveItems)
			g.GET("attribute/items", MustLogin(), api.Attribute.GetItems)
		}
		{
			g.GET("product/list", MustLogin(), api.Product.GetList)
			g.GET("product/detail", MustLogin(), api.Product.GetByID)
			g.POST("product/save", MustLogin(), api.Product.Save)
			g.POST("product/delete", MustLogin(), api.Product.DelByID)
		}
		{
			g.GET("order/list", MustLogin(), api.Order.GetList)
			g.GET("order/detail", MustLogin(), api.Order.GetByOrderID)
			g.POST("order/delivery", MustLogin(), api.Order.Delivery)
			g.GET("order/plan/list", MustLogin(), api.Order.GetNeedSendList)
			g.POST("order/cancel", MustLogin(), api.Order.CancelOrder)
			g.POST("order/update/address", MustLogin(), api.Order.UpdateAddress)
			g.GET("order/query/express", MustLogin(), api.Order.QueryExpress)
		}
		{
			g.GET("banner/detail", MustLogin(), api.Banner.GetByID)
			g.GET("banner/list", MustLogin(), api.Banner.GetList)
			g.POST("banner/save", MustLogin(), api.Banner.Save)
			g.POST("banner/delete", MustLogin(), api.Banner.Del)
		}
		{
			g.GET("district/list", MustLogin(), api.District.GetGradation)
			g.GET("freight/list", MustLogin(), api.Freight.GetList)
			g.POST("freight/save", MustLogin(), api.Freight.SaveAll)
		}
		{
			g.GET("statistics/statistics", MustLogin(), api.Statistics.Statistics)
		}
	}
}
