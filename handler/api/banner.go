package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/models"
)

type banner int

type BannerDefault banner

func (banner) GetByID(c *gin.Context) {
	var args struct {
		ID int `json:"id" form:"id"`
	}
	if c.Bind(&args) != nil {
		return
	}
	if data, err := models.BannerDefault.GetByID(args.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, data)
	}
}

func (banner) GetList(c *gin.Context) {
	if list, err := models.BannerDefault.GetList(0); err != nil {
		JSON(c, err)
	} else {
		JSON(c, list)
	}
}

func (banner) Save(c *gin.Context) {
	var args models.Banner
	if c.Bind(&args) != nil {
		return
	}
	var banner *models.Banner
	var err error
	if args.ID != 0 {
		if banner, err = models.BannerDefault.GetByID(args.ID); err != nil {
			JSON(c, err)
		}
	}
	if banner == nil {
		banner = new(models.Banner)
		banner.IsDel = models.SureNo
	}
	banner.Title = args.Title
	banner.Position = args.Position
	banner.Link = args.Link
	banner.Image = args.Image
	banner.Sort = args.Sort
	if err := banner.Save(); err != nil {
		JSON(c, err)
	} else {
		JSON(c, banner)
	}
}

func (banner) Del(c *gin.Context) {
	var args struct {
		ID int `json:"id"`
	}
	if err := models.BannerDefault.DelByID(args.ID); err != nil {
		JSON(c, err)
	} else {
		JSON(c, "success")
	}
}
