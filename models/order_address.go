package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OrderAddress struct {
	Model
	OrderId       string `json:"order_id"`
	UserID        int    `json:"user_id"`
	AddressID     int    `json:"address_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	DistrictCode  string `json:"district_code"`
	DistrictName  string `json:"district_name"`
	DetailAddress string `json:"detail_address"`
	Created       int64  `json:"created"`
}

func (addr *OrderAddress) Insert(db *gorm.DB) error {
	addr.Created = time.Now().Unix()
	return db.Create(addr).Error
}
