package modules

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
)

type adminActionLog int

var AdminActionLog adminActionLog

// 添加管理员操作日志
func (a adminActionLog) Insert(adminID int, c *gin.Context) {
	Remark := map[string]map[string]string{
		"/api/admin_user/group":      {"GET": "查看管理员权限组列表"},
		"/api/admin_user/group/info": {"GET": "查看管理员权限组"},
		"/api/admin_user/group/save": {"POST": "保存管理员权限组"},
		"/api/admin_user/group/del":  {"POST": "删除管理员权限组"},
		"/api/admin_user/users":      {"GET": "查看管理员账号列表"},
		"/api/admin_user/user/info":  {"GET": "编辑管理员账号"},
		"/api/admin_user/user/save":  {"POST": "保存管理员账号"},
		"/api/admin_user/user/del":   {"POST": "删除管理员账号"},
		"/api/admin_user/role/list":  {"GET": "获取权限列表"},

		"/api/passport/update/password": {"GET": "修改密码"},
		"/api/passport/login":           {"POST": "登录"},
		"/api/passport/login/out":       {"POST": "退出"},
		"/api/passport/user":            {"GET": "获取管理员登录信息"},

		"/api/app_version/list": {"GET": "app版本吧"},
		"/api/app_version/save": {"POST": "保存app版本吧"},

		"/api/article/list": {"GET": "查看新闻列表"},
		"/api/article/save": {"POST": "保存新闻"},
		"/api/article/del":  {"POST": "删除新闻"},

		"/api/product/list": {"GET": "查看产品列表"},
		"/api/product/save": {"POST": "保存产品"},
		"/api/product/del":  {"POST": "删除产品"},
	}
	bodyJson, _ := json.Marshal(c.Request.Body)
	RequestURI := string(c.Request.URL.Path)
	remark := "非法操作"
	if v, ok := Remark[RequestURI]; ok {
		if val, ok := v[c.Request.Method]; ok {
			remark = val
		}
	}
	logModels := models.AdminActionLog{}
	logModels.Ip = c.ClientIP()
	logModels.AdminID = adminID
	logModels.Path = c.Request.Host + c.Request.RequestURI
	logModels.Body = string(bodyJson)
	logModels.Remark = remark
	logModels.Created = time.Now().Unix()
	logModels.Insert()
}
