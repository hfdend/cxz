package models

import (
	"time"
)

// A User 用户
// swagger:model User
type User struct {
	Model
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
}

// TableName TableName
func (User) TableName() string {
	return "user"
}

// Insert Insert
func (u *User) Insert() error {
	u.Created = time.Now().Unix()
	return u.DB().Create(u).Error
}
