package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
)

type product int

var Product product

func (product) GetList(c *gin.Context) {
	var args struct {
		Page int `json:"page" form:"page"`
	}
	if c.Bind(&args) != nil {
		return
	}
	pager := models.NewPager(args.Page, 20)
	if list, err := models.ProductDefault.GetList(pager); err != nil {
		JSON(c, err)
	} else {
		JSON(c, map[string]interface{}{
			"list":  list,
			"pager": pager,
		})
	}
}

func (product) GetByID(c *gin.Context) {
	var args struct {
		ID int `json:"id" form:"id"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if product, err := models.ProductDefault.GetByID(args.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, product)
	}
}

func (product) Save(c *gin.Context) {
	var args models.Product
	if c.Bind(&args) != nil {
		return
	}
	var product *models.Product
	var err error
	if args.ID != 0 {
		if product, err = models.ProductDefault.GetByID(args.ID); err != nil {
			JSON(c, err)
			return
		} else if product == nil {
			JSON(c, errors.New("商品不存在"))
			return
		}
	} else {
		product = new(models.Product)
		product.IsDel = models.SureNo
	}
	product.Name = args.Name
	product.Type = args.Type
	product.Taste = args.Taste
	product.Unit = args.Unit
	product.Price = args.Price
	product.Image = args.Image
	product.Mark = args.Mark
	product.Intro = args.Intro
	if err = product.Save(); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}

func (product) DelByID(c *gin.Context) {
	var args struct {
		ID int `json:"id" form:"id"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if err := models.ProductDefault.DelByID(args.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}
