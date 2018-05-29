package models

import (
	"time"

	"github.com/hfdend/cxz/cli"
)

type AttributeItem struct {
	Model
	AttributeID int    `json:"attribute_id"`
	Sort        int    `json:"sort"`
	Name        string `json:"name"`
	Created     int64  `json:"created"`
}

var AttributeItemDefault AttributeItem

func (AttributeItem) TableName() string {
	return "attribute_item"
}

func (item *AttributeItem) Insert() error {
	if item.Created == 0 {
		item.Created = time.Now().Unix()
	}
	return item.DB().Create(item).Error
}

func (item *AttributeItem) DelByAttributeID(id int) error {
	return item.DB().Delete(AttributeItem{}, "attribute_id = ?", id).Error
}

func (AttributeItem) GetByAttributeID(id int) (list []*AttributeItem, err error) {
	err = cli.DB.Where("attribute_id = ?", id).Order("sort asc").Find(&list).Error
	return
}
