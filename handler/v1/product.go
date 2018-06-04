package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/modules"
)

type product int

var Product product

// 商品列表
// swagger:response ProductGetListResp
type ProductGetListResp struct {
	// in: body
	Body struct {
		List  []*models.Product `json:"list"`
		Pager *models.Pager
	}
}

// swagger:parameters Product_GetList
type ProductGetListArgs struct {
	Page int `json:"page" form:"page"`
	models.ProductCondition
}

// swagger:route GET /product/list 商品 Product_GetList
// 获取商品列表
// responses:
//     200: ProductGetListResp
func (product) GetList(c *gin.Context) {
	var args ProductGetListArgs
	var resp ProductGetListResp
	var err error
	if c.Bind(&args) != nil {
		return
	}
	fmt.Println(args)
	if resp.Body.List, err = modules.Product.GetList(args.ProductCondition, resp.Body.Pager); err != nil {
		JSON(c, err)
	} else {
		JSON(c, resp.Body)
	}
}

// 商品分类和属性
// swagger:response ProductAttributeItemsResp
type ProductAttributeItemsResp struct {
	// in: body
	Body []*models.AttributeItem
}

// swagger:route GET /product/attribute/items 商品 Product_AttributeItems
// responses:
//     200: ProductAttributeItemsResp
func (product) AttributeItems(c *gin.Context) {
	list, err := models.AttributeItemDefault.GetAll()
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}
