package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// A User 用户
// swagger:model User
type User struct {
	Model
	Phone    string `json:"phone"`
	Password string `json:"-"`
	UnionID  string `json:"unionid"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
}

var UserDefault User

// TableName TableName
func (User) TableName() string {
	return "user"
}

// Insert Insert
func (u *User) Insert() (int64, error) {
	u.Created = time.Now().Unix()
	u.Updated = time.Now().Unix()
	return DBInsertIgnore(u.DB(), u)
}

func (u User) GetByID(id int) (data *User, err error) {
	data = new(User)
	err = u.DB().Where("id = ?", id).Find(data).Error
	return
}

func (u User) GetByPhone(phone string) (data *User, err error) {
	data = new(User)
	if err = u.DB().Where("phone = ?", phone).Find(data).Error; gorm.IsRecordNotFoundError(err) {
		err = nil
	}
	return
}

func (u User) GetByUnionID(unionID string) (data *User, err error) {
	data = new(User)
	if err = u.DB().Where("unionid = ?", unionID).Find(data).Error; gorm.IsRecordNotFoundError(err) {
		err = nil
	}
	return
}
