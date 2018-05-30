package models

import "github.com/hfdend/cxz/cli"

type Attribute struct {
	Model
	Name  string           `json:"name"`
	Kind  string           `json:"kind"`
	Items []*AttributeItem `json:"items" gorm:"-"`
}

var AttributeDefault Attribute

func (Attribute) TableName() string {
	return "attribute"
}

func (attr Attribute) GetList() (list []*Attribute, err error) {
	err = attr.DB().Find(&list).Error
	return
}

func (Attribute) GetAll() (list []*Attribute, err error) {
	if err = cli.DB.Find(&list).Error; err != nil {
		return
	}
	var items []*AttributeItem
	if items, err = AttributeItemDefault.GetAll(); err != nil {
		return
	}
	attrMap := map[int]*Attribute{}
	for _, v := range list {
		attrMap[v.ID] = v
	}
	for _, v := range items {
		if attr, ok := attrMap[v.AttributeID]; ok {
			attr.Items = append(attr.Items, v)
		}
	}
	return
}
