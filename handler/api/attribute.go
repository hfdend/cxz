package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
)

type attribute int

var Attribute attribute

func (attribute) GetList(c *gin.Context) {
	list, err := models.AttributeDefault.GetList()
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

func (attribute) GetAll(c *gin.Context) {
	list, err := models.AttributeDefault.GetAll()
	if err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

func (attribute) SaveItems(c *gin.Context) {
	var args struct {
		ID    int                     `json:"id"`
		Items []*models.AttributeItem `json:"items"`
	}
	if err := c.Bind(&args); err != nil {
		return
	}
	if err := models.AttributeItemDefault.DelByAttributeID(args.ID); err != nil {
		JSON(c, err)
		return
	}
	for _, v := range args.Items {
		v.AttributeID = args.ID
		if err := v.Insert(); err != nil {
			JSON(c, err)
			return
		}
	}
	JSON(c, "success")
}

func (attribute) GetItems(c *gin.Context) {
	var args struct {
		ID int `json:"id" form:"id"`
	}
	if err := c.Bind(&args); err != nil {
		return
	}
	if list, err := models.AttributeItemDefault.GetByAttributeID(args.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}
