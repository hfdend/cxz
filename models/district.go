package models

import (
	"sync"
)

type District struct {
	Model
	Name     string      `json:"name"`
	ParentID int         `json:"parent_id"`
	Initial  string      `json:"initial"`
	Initials string      `json:"initials"`
	Pinyin   string      `json:"pinyin"`
	Extra    string      `json:"extra"`
	Suffix   string      `json:"suffix"`
	Code     string      `json:"code"`
	AreaCode string      `json:"area_code"`
	Order    int         `json:"order"`
	Children []*District `json:"children" gorm:"-"`
}

var (
	districtCodeMapping = map[string]*District{}
	districtIDMapping   = map[int]*District{}
	districtGradation   []*District
	districtOnce        sync.Once
)

var DistrictDefault District

func (District) TableName() string {
	return "district"
}

func (d District) once() error {
	var err error
	districtOnce.Do(func() {
		var list []*District
		if err = d.DB().Find(&list).Error; err != nil {
			return
		}
		for _, v := range list {
			districtCodeMapping[v.Code] = v
			districtIDMapping[v.ID] = v
		}
		for _, v := range list {
			if v.ParentID == 0 {
				districtGradation = append(districtGradation, v)
			} else if p, ok := districtIDMapping[v.ParentID]; ok {
				p.Children = append(p.Children, v)
			}
		}
	})
	return err
}

func (d District) GetGradation() (list []*District, err error) {
	if err = d.once(); err != nil {
		return
	}
	list = districtGradation
	return
}

func (d District) GetNames(code string) (list []string, err error) {
	if err = d.once(); err != nil {
		return
	}
	data, ok := districtCodeMapping[code]
	if !ok {
		return
	}
	list = append(list, data.Name)
	id := data.ParentID
	for {
		data, ok := districtIDMapping[id]
		if !ok {
			break
		}
		list = append([]string{data.Name}, list...)
		id = data.ParentID
	}
	return
}
