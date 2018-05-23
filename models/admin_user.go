package models

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

const AdminUserTokenKey = "admin_user_token_"

type AdminUser struct {
	Model
	Username string `json:"username"`
	Password string `json:"password"`
	GroupId  int    `json:"group_id"`
	Created  int64  `json:"created"`

	AdminGroup *AdminGroup `gorm:"-" json:"admin_group"`
}

var AdminUserDefault AdminUser

func (AdminUser) TableName() string {
	return "admin_user"
}

func (a *AdminUser) Save() error {
	if a.Created == 0 {
		a.Created = time.Now().Unix()
	}
	return a.DB().Save(a).Error
}

func (a AdminUser) GetByUsername(username string) (*AdminUser, error) {
	var res AdminUser
	if err := a.DB().Where("username = ?", username).Where("is_del = ?", 0).Find(&res).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a AdminUser) GetById(id int) (*AdminUser, error) {
	var res AdminUser
	if err := a.DB().Where("id = ?", id).Where("is_del = ?", 0).Find(&res).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a AdminUser) GetToken(id int) (string, error) {
	token := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
	key := fmt.Sprintf("%s%s", AdminUserTokenKey, token)
	err := cli.Redis.Set(key, id, 20*24*time.Hour).Err()
	return token, err
}

func (a AdminUser) GetUserByToken(token string) (*AdminUser, error) {
	var (
		s   string
		id  int
		err error
	)
	key := fmt.Sprintf("%s%s", AdminUserTokenKey, token)
	if s, err = cli.Redis.Get(key).Result(); err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else if id, err = strconv.Atoi(s); err != nil {
		return nil, err
	}
	return a.GetByIdWithGroup(id)
}

func (a AdminUser) DelToken(token string) error {
	key := fmt.Sprintf("%s%s", AdminUserTokenKey, token)
	return cli.Redis.Del(key).Err()
}

func (a AdminUser) GetByGroupId(id int) (list []*AdminUser, err error) {
	err = a.DB().Where("group_id = ?", id).Find(&list).Error
	return
}

func (a AdminUser) GetList() (list []*AdminUser, err error) {
	err = a.DB().Where("is_del = ?", 0).Find(&list).Error
	return
}

func (a AdminUser) DelById(id int) error {
	return a.DB().Delete(a, "id = ?", id).Error
}

func (a *AdminUser) SetGroup() error {
	gp, err := AdminGroupDefault.GetById(a.GroupId)
	if err != nil {
		return err
	}
	a.AdminGroup = gp
	return nil
}

func (a AdminUser) GetByIdWithGroup(id int) (*AdminUser, error) {
	u, err := a.GetById(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, err
	}
	if err := u.SetGroup(); err != nil {
		return nil, err
	}
	return u, nil
}

func (a AdminUser) DelUser(id int) (err error) {
	return a.DB().Table(a.TableName()).Where("id = ?", id).Update("is_del", 1).Error
}
