package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxr/errors"
	"github.com/hfdend/cxr/models"
	"github.com/hfdend/cxr/modules"
)

type passport int

var Passport passport

func (passport) RegisterSend(c *gin.Context) {
	var args struct {
		Phone string `json:"phone"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if _, err := modules.Passport.SendRegisterCode(args.Phone); err != nil {
		JSON(c, err)
	} else {
		JSON(c, SUCCESS)
	}
}

func (passport) Register(c *gin.Context) {
	var args struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if _, err := modules.Passport.Register(args.Phone, args.Code, args.Password); err != nil {
		JSON(c, err)
	} else {
		JSON(c, SUCCESS)
	}
}

func (passport) Login(c *gin.Context) {
	var args struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if token, err := modules.Passport.Login(args.Phone, args.Password); err != nil {
		JSON(c, err)
	} else {
		JSON(c, token)
	}
}

func SetLoginUser(c *gin.Context) {
	accessToken := c.Request.Header.Get("token")
	if accessToken == "" {
		if cookie, err := c.Request.Cookie("token"); err != nil {
			JSON(c, err)
			return
		} else if cookie != nil {
			accessToken = cookie.Value
		}
	}
	if accessToken == "" {
		return
	}
	userId, err := models.TokenDefault.GetUserId(accessToken)
	if err != nil {
		JSON(c, err)
		return
	} else if userId == 0 {
		return
	}
	user, err := models.UserDefault.GetByID(userId)
	if err != nil {
		JSON(c, err)
		return
	} else if user == nil {
		return
	}
	c.Set("_login_user", user)
}

func GetUser(c *gin.Context) *models.User {
	if v, ok := c.Get("_login_user"); ok {
		return v.(*models.User)
	}
	return nil
}

func MustLogin(c *gin.Context) {
	if user := GetUser(c); user == nil {
		JSON(c, errors.New("请登录", errors.NoLogin))
	}
}
