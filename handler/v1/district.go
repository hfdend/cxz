package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
)

type district int

var District district

// 地区数据
// swagger:response DistrictGetGradationResp
type DistrictGetGradationResp struct {
	// in: body
	Body []*models.District
}

// swagger:route GET /district/gradation 地址 District_GetGradation
// 获取地区数据
// responses:
//     200: DistrictGetGradationResp
func (district) GetGradation(c *gin.Context) {
	list, err := models.DistrictDefault.GetGradation()
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}
