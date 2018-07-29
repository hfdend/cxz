package modules

import (
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/utils"
)

type adminUser int

var AdminUser adminUser

func (adminUser) GetUserList() (list []*models.AdminUser, err error) {
	list, err = models.AdminUserDefault.GetList()
	if err != nil {
		return
	}
	for _, v := range list {
		if err = v.SetGroup(); err != nil {
			return
		}
	}
	return
}

func (adminUser) DelUser(id int) error {
	return models.AdminUserDefault.DelUser(id)
}

func (adminUser) SaveUser(user *models.AdminUser) error {
	if user.Username == "" {
		return errors.New("请填写账号")
	}
	var oldUser *models.AdminUser
	var err error
	if oldUser, err = models.AdminUserDefault.GetByUsername(user.Username); err != nil {
		return err
	} else if oldUser != nil && oldUser.ID != user.ID {
		return errors.New("账号已存在")
	}
	if user.Password != "" {
		user.Password = utils.EncodePassword(user.Password)
	} else {
		user.Password = oldUser.Password
	}
	if user.ID == 0 {
		if user.Password == "" {
			return errors.New("请填写密码")
		}
	}
	if g, err := models.AdminGroupDefault.GetByID(user.GroupID); err != nil {
		return err
	} else if g == nil {
		return errors.New("用户组不存在")
	}
	if err := user.Save(); err != nil {
		return err
	}
	return nil
}

func (adminUser) GetByUsername(username string) (*models.AdminUser, error) {
	return models.AdminUserDefault.GetByUsername(username)
}

func (a adminUser) GetByID(id int) (*models.AdminUser, error) {
	return models.AdminUserDefault.GetByID(id)
}

func (a adminUser) GetToken(id int) (string, error) {
	return models.AdminUserDefault.GetToken(id)
}

func (a adminUser) GetUserByToken(token string) (*models.AdminUser, error) {
	return models.AdminUserDefault.GetUserByToken(token)
}

func (a adminUser) DelToken(token string) error {
	return models.AdminUserDefault.DelToken(token)
}

func (a adminUser) GetGroup() ([]*models.AdminGroup, error) {
	return models.AdminGroupDefault.GetList()
}

func (a adminUser) GetGroupByID(id int) (*models.AdminGroup, error) {
	return models.AdminGroupDefault.GetByID(id)
}

func (a adminUser) SaveGroup(group *models.AdminGroup) error {
	if group.Name == "" {
		return errors.New("请填写组名")
	}
	if err := group.Save(); err != nil {
		return err
	}
	return nil
}

func (a adminUser) DelGroup(id int) error {
	if list, err := models.AdminUserDefault.GetByGroupID(id); err != nil {
		return err
	} else if len(list) > 0 {
		return errors.New("分组下存在用户不能删除")
	}
	if err := models.AdminGroupDefault.DelByID(id); err != nil {
		return err
	}
	return nil
}
