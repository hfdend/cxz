package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/jinzhu/gorm"
)

type BannerPosition int

const (
	BannerPositionAll BannerPosition = iota
)

type Banner struct {
	Model
	Title    string         `json:"title"`
	Position BannerPosition `json:"position"`
	Link     string         `json:"link"`
	Image    string         `json:"image"`
	IsDel    Sure           `json:"is_del"`
	Sort     int            `json:"sort"`
	Created  int64          `json:"created"`
	Updated  int64          `json:"updated"`
	DelTime  int64          `json:"del_time"`
	ImageSrc string         `json:"image_src" gorm:"-"`
}

var BannerDefault Banner

func (Banner) TableName() string {
	return "banner"
}

func (bn *Banner) Save() error {
	if bn.Created == 0 {
		bn.Created = time.Now().Unix()
	}
	bn.Updated = time.Now().Unix()
	return cli.DB.Save(bn).Error
}

func (Banner) GetList(pos int) (list []*Banner, err error) {
	db := cli.DB.Model(Banner{})
	if pos != 0 {
		db = db.Where("position = ? and is_del = ?", pos, SureNo)
	}
	if err = db.Order("sort desc").Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		v.SetImageSrc()
	}
	return
}

func (Banner) GetByID(id int) (*Banner, error) {
	var data Banner
	if err := cli.DB.Where("id = ? and is_del = ?", id, SureNo).Find(&data).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	data.SetImageSrc()
	return &data, nil
}

func (Banner) DelByID(id int) error {
	data := map[string]interface{}{
		"del_time": time.Now().Unix(),
		"is_del":   SureYes,
	}
	return cli.DB.Model(Banner{}).Where("id = ?", id).Update(data).Error
}

func (bn *Banner) SetImageSrc() {
	if bn.Image == "" {
		return
	}
	c := conf.Config.Aliyun.OSS
	bn.ImageSrc = fmt.Sprintf("%s/%s", strings.TrimRight(c.Domain, "/"), strings.TrimLeft(bn.Image, "/"))
	return
}
