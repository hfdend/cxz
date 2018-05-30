package api

import (
	"fmt"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/modules"
)

type file int

var File file

func (file) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		JSON(c, err)
		return
	}
	extName := strings.ToLower(path.Ext(header.Filename))
	switch extName {
	default:
		JSON(c, errors.New("不允许的上传格式"))
		return
	case ".png", ".jpg", ".jpeg", ".gif":
	}
	uid := uuid.New().String()
	key := fmt.Sprintf("admin/%s%s", strings.Replace(uid, "-", "/", -1), extName)
	cnd, err := modules.File.UploadToCDN(file, key)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, cnd)
	}
}
