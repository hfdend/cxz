package models

type AdminGroup struct {
	Model
	Name  string `json:"name"`
	Roles string `json:"roles"`
}

var AdminGroupDefault AdminGroup

func (AdminGroup) TableName() string {
	return "admin_group"
}

func (g AdminGroup) GetById(id int) (*AdminGroup, error) {
	var res AdminGroup
	if db := g.DB().Table(g.TableName()).Where("id = ?", id).Scan(&res); db.RecordNotFound() {
		return nil, nil
	} else if db.Error != nil {
		return nil, db.Error
	}
	return &res, nil
}

func (g AdminGroup) GetList() (list []*AdminGroup, err error) {
	err = g.DB().Find(&list).Error
	return
}

func (g *AdminGroup) Insert() error {
	return g.DB().Create(g).Error
}

func (g *AdminGroup) Save() error {
	return g.DB().Save(g).Error
}

func (g AdminGroup) DelById(id int) error {
	return g.DB().Delete(g, "id = ?", id).Error
}
