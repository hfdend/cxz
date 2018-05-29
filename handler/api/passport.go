package api

import (
	"net/http"

	"gitee.com/cardctl/server/role"
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
	"github.com/hfdend/cxz/utils"
)

type passport int

var Passport passport

func (passport) Login(c *gin.Context) {
	var args struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if c.Bind(&args) != nil {
		return
	}
	var adminUser *models.AdminUser
	var err error
	var token string
	if adminUser, err = modules.AdminUser.GetByUsername(args.Username); err != nil {
		JSON(c, err)
	} else if adminUser == nil {
		JSON(c, errors.New("账号不存在"))
	} else if utils.EncodePassword(args.Password) != adminUser.Password {
		JSON(c, errors.New("密码错误"))
	} else if token, err = modules.AdminUser.GetToken(adminUser.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, token)
	}
}

func (passport) LoginOut(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if err := modules.AdminUser.DelToken(token); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (passport) GetUser(c *gin.Context) {
	user := getUser(c)
	if user.ID != 0 {
		JSON(c, user)
	} else {
		JSON(c, nil)
	}
}

func (passport) UpdatePassword(c *gin.Context) {
	var args struct {
		OldPassword string
		NewPassword string
	}
	if err := c.Bind(&args); err != nil {
		return
	}
	user := getUser(c)
	if utils.EncodePassword(args.OldPassword) != user.Password {
		JSON(c, errors.New("旧密码错误"))
		return
	}
	user.Password = utils.EncodePassword(args.NewPassword)
	if err := user.Save(); err != nil {
		JSON(c, errors.New("旧密码错误"))
		return
	}
	JSON(c, "success")
}

func (passport) SetLoginUser(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		ck, err := c.Request.Cookie("token")
		if err == http.ErrNoCookie {
			return
		}
		if err != nil {
			JSON(c, err)
			return
		}
		token = ck.Value
	}
	if token == "" {
		return
	}
	if user, err := modules.AdminUser.GetUserByToken(token); err != nil {
		JSON(c, err)
	} else {
		setUser(c, user)
	}
}

func (passport) MustLogin(or ...[]string) gin.HandlerFunc {
	var fn gin.HandlerFunc
	fn = func(c *gin.Context) {
		user := getUser(c)
		if user == nil {
			JSON(c, errors.New("请登录", errors.NoLogin))
			return
		}
		if err := role.Check(user.AdminGroup.Roles, or...); err == role.ErrNoRole {
			JSON(c, errors.New("权限不足"))
			return
		} else if err != nil {
			JSON(c, err)
			return
		}
	}
	return fn
}

func setUser(c *gin.Context, user *models.AdminUser) {
	c.Set("_current_user", user)
}

func getUser(c *gin.Context) *models.AdminUser {
	if v, ok := c.Get("_current_user"); ok {
		return v.(*models.AdminUser)
	}
	return nil
}
