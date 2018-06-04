package v1

import (
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
	resp.Body.Pager = models.NewPager(args.Page, 20)
	if resp.Body.List, err = modules.Product.GetList(args.ProductCondition, resp.Body.Pager); err != nil {
		JSON(c, err)
	} else {
		JSON(c, resp.Body)
	}
}

// swagger:parameters Product_AttributeItems
type ProductAttributeItemsArgs struct {
	// 1: 获取商品分类选项
	// 2: 获取口味选项
	Id int `json:"id" form:"id"`
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
	var args ProductAttributeItemsArgs
	if c.Bind(&args) != nil {
		return
	}
	list, err := models.AttributeItemDefault.GetByAttributeID(args.Id)
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}
