package v1

import (
	"github.com/gin-gonic/gin"
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
}
