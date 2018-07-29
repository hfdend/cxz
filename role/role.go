package role

import (
	"errors"
	"strings"
)

type Role struct {
	Filed string `json:"filed"`
	Name  string `json:"name"`
}

type Group struct {
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

var GroupList []Group
var ErrNoRole = errors.New("no role")

const (
	Administrator      = "administrator"
	AdminUser          = "admin_user"
	AdminUserEdit      = "admin_user_edit"
	AdminUserDel       = "admin_user_del"
	AdminUserGroup     = "admin_user_group"
	AdminUserGroupInfo = "admin_user_group_info"
	AdminUserGroupEdit = "admin_user_group_edit"
	AdminUserGroupDel  = "admin_user_group_del"
	AdminUserRole      = "admin_user_role"
)

func init() {
	{
		g := Group{Name: "管理员"}
		g.Roles = append(g.Roles, Role{Administrator, "超级管理员"})
		g.Roles = append(g.Roles, Role{AdminUser, "账号列表"})
		g.Roles = append(g.Roles, Role{AdminUserEdit, "编辑账号"})
		g.Roles = append(g.Roles, Role{AdminUserDel, "删除账号"})
		g.Roles = append(g.Roles, Role{AdminUserGroup, "权限组列表"})
		g.Roles = append(g.Roles, Role{AdminUserGroupInfo, "查看权限组"})
		g.Roles = append(g.Roles, Role{AdminUserGroupEdit, "编辑权限组"})
		g.Roles = append(g.Roles, Role{AdminUserGroupDel, "删除权限组"})
		g.Roles = append(g.Roles, Role{AdminUserRole, "获取管理员权限"})
		GroupList = append(GroupList, g)
	}
	{
		//g := Group{Name: "设置"}
	}
}

func Check(userRoles string, or ...[]string) error {
	userRolesMap := map[string]bool{}
	for _, v := range strings.Split(userRoles, ",") {
		userRolesMap[v] = true
	}
	if _, ok := userRolesMap[Administrator]; ok {
		return nil
	}
	if len(or) == 0 {
		return nil
	}
	for _, ary := range or {
		yes := true
		for _, v := range ary {
			if _, ok := userRolesMap[v]; !ok {
				yes = false
				break
			}
		}
		if yes {
			return nil
		}
	}
	return ErrNoRole
}
