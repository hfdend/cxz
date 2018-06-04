package modules

import "github.com/hfdend/cxz/models"

type product int

var Product product

func (product) GetList(cond models.ProductCondition, pager *models.Pager) (list []*models.Product, err error) {
	list, err = models.ProductDefault.GetList(cond, pager)
	return
}
