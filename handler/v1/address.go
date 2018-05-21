package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
)

type address int

var Address address

// swagger:parameters Address_Save
type AddressSaveArgs struct {
	// in: body
	Body struct {
		// 有修改，没有新增
		ID int `json:"id"`
		// 收货人姓名
		Name string `json:"name"`
		// 收货人电话
		Phone string `json:"phone"`
		// 地区code
		DistrictCode string `json:"district_code"`
		// 详细地址
		DetailAddress string `json:"detail_address"`
	}
}

// 地址数据
// swagger:response AddressSaveResp
type AddressSaveResp struct {
	// in: body
	Body *models.Address
}

// swagger:route POST /address/save 地址 Address_Save
// 保存地址
// responses:
//     200: AddressSaveResp
func (address) Save(c *gin.Context) {
	var args AddressSaveArgs
	if c.Bind(&args.Body) != nil {
		return
	}
	body := args.Body
	user := GetUser(c)
	if address, err := modules.Address.Save(body.ID, user.ID, body.Name, body.Phone, body.DistrictCode, body.DetailAddress); err != nil {
		JSON(c, err)
	} else {
		JSON(c, address)
	}
}
