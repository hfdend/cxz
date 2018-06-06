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
		DetailAddress string      `json:"detail_address"`
		IsDefault     models.Sure `json:"is_default"`
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
	if address, err := modules.Address.Save(body.ID, user.ID, body.Name, body.Phone, body.DistrictCode, body.DetailAddress, body.IsDefault); err != nil {
		JSON(c, err)
	} else {
		JSON(c, address)
	}
}

// swagger:parameters Address_Del
type AddressDelArgs struct {
	// in: body
	Body struct {
		ID int `json:"id"`
	}
}

// swagger:route POST /address/del 地址 Address_Del
// 删除地址
// responses:
//     200: SUCCESS
func (address) Del(c *gin.Context) {
	var args AddressDelArgs
	if c.Bind(&args.Body) != nil {
		return
	}
	user := GetUser(c)
	if err := models.AddressDefault.DelById(user.ID, args.Body.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, SUCCESS)
	}
}

// 地址列表
// swagger:response AddressListResp
type AddressListResp struct {
	// in: body
	Body []*models.Address
}

// swagger:route GET /address/list 地址 Address_List
// 获取地址列表
// responses:
//     200: AddressListResp
func (address) List(c *gin.Context) {
	user := GetUser(c)
	if list, err := models.AddressDefault.GetList(user.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

// swagger:parameters Address_GetByID
type AddressGetByIDArgs struct {
	ID int `json:"id" form:"id"`
}

// 地址详情
// swagger:response AddressGetByIDResp
type AddressGetByIDResp struct {
	// in: body
	Body *models.Address
}

// swagger:route GET /address/detail 地址 Address_GetByID
// 获取地址列表
// responses:
//     200: AddressGetByIDResp
func (address) GetByID(c *gin.Context) {
	var args AddressGetByIDArgs
	if c.Bind(&args) != nil {
		return
	}
	user := GetUser(c)
	if addr, err := models.AddressDefault.GetByIDAndUserID(args.ID, user.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, addr)
	}
}
