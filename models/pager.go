package models

import (
	"fmt"
	"math"

	"github.com/jinzhu/gorm"
)

type Pager struct {
	Limit     int `json:"limit"`
	Page      int `json:"page"`
	PageCount int `json:"page_count"`
	Offset    int `json:"offset"`
	Count     int `json:"count"`
}

func NewPager(currentPage, limit int) *Pager {
	if currentPage < 1 {
		currentPage = 1
	}
	p := new(Pager)
	p.Page = currentPage
	p.Limit = limit
	p.Offset = (currentPage - 1) * limit
	return p
}

func (p *Pager) GetLimit() string {
	return fmt.Sprintf("%v, %v", (p.Page-1)*p.Limit, p.Limit)
}

func (p *Pager) SetPager(count int) {
	p.Count = count
	p.PageCount = int(math.Ceil(float64(count) / float64(p.Limit)))
}

func (p *Pager) Exec(db *gorm.DB) (*gorm.DB, error) {
	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, err
	}
	p.SetPager(count)
	db = db.Limit(p.Limit).Offset(p.Offset)
	return db, nil
}
