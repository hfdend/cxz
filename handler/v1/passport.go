package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
)

type passport int

var Passport passport

// swagger:parameters Passport_RegisterSend
type PassportRegisterSendArgs struct {
	// in: body
	Body struct {
		Phone string `json:"phone"`
	}
}

// swagger:route POST /register/send 账号 Passport_RegisterSend
// 发送注册验证码
// responses:
//     200: SUCCESS
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

// swagger:parameters Passport_Register
type PassportRegisterArgs struct {
	// in: body
	Body struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}
}

// swagger:route POST /register 账号 Passport_Register
// 账号注册
// responses:
//     200: SUCCESS
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

// swagger:parameters Passport_Login
type PassportLoginArgs struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// swagger:route POST  /login 账号 Passport_Login
// 登录
// responses:
//     200: SUCCESS
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

// swagger:parameters Passport_LoginByJsCode
type PassportLoginByJsCodeArgs struct {
	// in: body
	Body struct {
		// 用户登录凭证
		JSCode string `json:"js_code"`
	}
}

// swagger:response PassportLoginByJsCodeResp
type PassportLoginByJsCodeResp struct {
	// in: body
	Body models.Token
}

// swagger:route POST /miniprogram/login 账号 Passport_LoginByJsCode
// 通过小程序登录
// responses:
//     200: PassportLoginByJsCodeResp
func (passport) LoginByJsCode(c *gin.Context) {
	var args struct {
		JSCode string `json:"js_code"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if token, err := modules.Passport.LoginByJsCode(args.JSCode); err != nil {
		JSON(c, err)
	} else {
		JSON(c, token)
	}
}

// swagger:parameters Passport_BindPhone
type PassportBindPhoneArgs struct {
	// in: body
	Body struct {
		Phone string `json:"phone"`
	}
}

// swagger:route POST /bind/phone 账号 Passport_BindPhone
// 绑定手机号
// responses:
//     200: SUCCESS
func (passport) BindPhone(c *gin.Context) {
	var args struct {
		Phone string `json:"phone"`
	}
	if c.Bind(&args) != nil {
		return
	}
	user := GetUser(c)
	if err := modules.Passport.BindPhone(user.ID, args.Phone); err != nil {
		JSON(c, err)
	} else {
		JSON(c, SUCCESS)
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
