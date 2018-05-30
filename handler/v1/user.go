package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
)

type user int

var User user

// 登录用户信息
// swagger:response UserGetUserInfoResp
type UserGetUserInfoResp struct {
	// in: body
	Body struct {
		*models.User
	}
}

// swagger:route GET /user 用户 User_GetUserInfo
// 获取登录用户信息
// responses:
//     200: UserGetUserInfoResp
func (user) GetUserInfo(c *gin.Context) {
	user := GetUser(c)
	if user != nil {
		if len(user.Phone) == 32 {
			user.Phone = ""
		}
	}
	JSON(c, user)
}
