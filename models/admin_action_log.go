package models

type AdminActionLog struct {
	Model
	Id      int    `gorm:"primary_key" json:"id"`
	AdminID int    `json:"admin_id"`
	Path    string `json:"path"`
	Remark  string `json:"remark"`
	Body    string `json:"body"`
	Ip      string `json:"ip"`
	Created int64  `json:"created"`
}

type AdminActionLogCondition struct {
	AdminId   int    `json:"admin_id"`
	Path      string `json:"path"`
	Ip        string `json:"ip"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

var AdminActionLogDefault AdminActionLog

func (AdminActionLog) TableName() string {
	return "admin_action_log"
}

func (a AdminActionLog) GetById(id int) (*AdminActionLog, error) {
	var res AdminActionLog
	if db := a.DB().Table(a.TableName()).Where("id = ?", id).Scan(&res); db.RecordNotFound() {
		return nil, nil
	} else if db.Error != nil {
		return nil, db.Error
	}
	return &res, nil
}

func (a *AdminActionLog) Insert() error {
	return a.DB().Create(a).Error
}

func (a AdminActionLog) DelById(id int) error {
	return a.DB().Delete(a, "id = ?", id).Error
}

func (a AdminActionLog) Search(condition AdminActionLogCondition, pager *Pager) (list []*User, err error) {
	db := a.DB().Table(a.TableName()).Order("id desc")
	if condition.AdminId != 0 {
		db = db.Where("admin_id = ?", condition.AdminId)
	}
	if condition.Ip != "" {
		db = db.Where("ip = ?", condition.Ip)
	}
	if condition.Path != "" {
		db = db.Where("path = ?", condition.Path)
	}
	if condition.StartTime != 0 {
		db = db.Where("created > ?", condition.StartTime)
	}
	if condition.EndTime != 0 {
		db = db.Where("created <= ?", condition.EndTime)
	}
	if pager != nil {
		if db, err = pager.Exec(db); err != nil {
			return
		} else if pager.Count == 0 {
			return
		}
	}
	err = db.Find(&list).Error
	return
}
