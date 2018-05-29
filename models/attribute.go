package models

type Attribute struct {
	Model
	Name string `json:"name"`
	Kind string `json:"kind"`
}

var AttributeDefault Attribute

func (Attribute) TableName() string {
	return "attribute"
}

func (attr Attribute) GetList() (list []*Attribute, err error) {
	err = attr.DB().Find(&list).Error
	return
}
