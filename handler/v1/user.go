package v1

import "github.com/gin-gonic/gin"

type user int

var User user

func (user) GetUserInfo(c *gin.Context) {
	user := GetUser(c)
	JSON(c, user)
}
