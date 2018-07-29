package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
	"github.com/hfdend/cxz/role"
)

type adminUser int

var AdminUser adminUser

func (adminUser) GetUserList(c *gin.Context) {
	if list, err := modules.AdminUser.GetUserList(); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

func (adminUser) GetUser(c *gin.Context) {
	var args struct {
		Id int
	}
	if c.Bind(&args) != nil {
		return
	}
	if res, err := models.AdminUserDefault.GetByID(args.Id); err != nil {
		JSON(c, err)
	} else {
		JSON(c, res)
	}
}

func (adminUser) SaveUser(c *gin.Context) {
	var args models.AdminUser
	if c.Bind(&args) != nil {
		return
	}
	if err := modules.AdminUser.SaveUser(&args); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (adminUser) DelUser(c *gin.Context) {
	var args struct {
		Id int
	}
	if c.Bind(&args) != nil {
		return
	}
	if err := modules.AdminUser.DelUser(args.Id); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (adminUser) GetGroupList(c *gin.Context) {
	if list, err := modules.AdminUser.GetGroup(); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

func (adminUser) GetGroupById(c *gin.Context) {
	var args struct {
		Id int
	}
	if c.Bind(&args) != nil {
		return
	}
	if res, err := modules.AdminUser.GetGroupByID(args.Id); err != nil {
		JSON(c, err)
	} else {
		JSON(c, res)
	}
}

func (adminUser) GetGroupByID(c *gin.Context) {
	var args struct {
		Id int
	}
	if err := c.Bind(&args); err != nil {
		JSON(c, err)
		return
	}
	if res, err := modules.AdminUser.GetGroupByID(args.Id); err != nil {
		JSON(c, err)
	} else {
		JSON(c, res)
	}
}

func (adminUser) GroupSave(c *gin.Context) {
	var args models.AdminGroup
	if c.Bind(&args) != nil {
		return
	}
	if err := modules.AdminUser.SaveGroup(&args); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (adminUser) GroupDel(c *gin.Context) {
	var args struct {
		Id int
	}
	if c.Bind(&args) != nil {
		return
	}
	if err := modules.AdminUser.DelGroup(args.Id); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (adminUser) RoleList(c *gin.Context) {
	JSON(c, role.GroupList)
}
