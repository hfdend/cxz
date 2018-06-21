package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
)

type banner int

var Banner banner

type BannerGetListArgs struct {
	// 某一个位置的Banner
	// 不传查询全部
	// 1 按月订购-狗
	// 2 按月订购-猫
	// 3 自主拼选-狗
	// 4 自主拼选-猫
	Position string `json:"position"`
}

// Banner列表
// swagger:response BannerGetListResp
type BannerGetListResp struct {
	// in: body
	Body []*models.Banner
}

// swagger:route GET /banner/list Banner Banner_GetList
// 获取Banner列表
// responses:
//     200: BannerGetListResp
func (banner) GetList(c *gin.Context) {
	var args BannerGetListArgs
	if c.Bind(&args) != nil {
		return
	}
	if list, err := models.BannerDefault.GetList(args.Position); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}
