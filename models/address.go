package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 用户地址
// swagger:model Address
type Address struct {
	Model
	UserID        int    `json:"user_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	DistrictCode  string `json:"district_code"`
	DistrictName  string `json:"district_name"`
	DetailAddress string `json:"detail_address"`
	IsDefault     Sure   `json:"is_default"`
	IsDel         Sure   `json:"is_del"`
	Created       int64  `json:"created"`
	Updated       int64  `json:"updated"`
}

var AddressDefault Address

func (Address) TableName() string {
	return "address"
}

func (a *Address) Save() error {
	a.Updated = time.Now().Unix()
	if a.Created == 0 {
		a.Created = time.Now().Unix()
	}
	return a.DB().Save(a).Error
}

func (a Address) UpdateNoDefault(userID, defaultID int) error {
	data := map[string]interface{}{
		"is_default": SureNo,
		"updated":    time.Now().Unix(),
	}
	return a.DB().Where("user_id = ? and id != ?", userID, defaultID).Update(data).Error
}

func (a Address) GetByID(id int) (*Address, error) {
	var data Address
	if err := a.DB().Where("id = ? and is_del = ?", id, SureNo).Find(&data).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &data, nil
}

func (a Address) DelById(userID, id int) error {
	data := map[string]interface{}{
		"is_del":  SureYes,
		"updated": time.Now().Unix(),
	}
	return a.DB().Table(a.TableName()).Where("id = ? and user_id = ?", id, userID).Update(data).Error
}

func (a Address) GetList(userID int) (list []*Address, err error) {
	err = a.DB().Where("user_id = ?", userID).Find(&list).Error
	return
}
