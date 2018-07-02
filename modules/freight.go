package modules

import "github.com/hfdend/cxz/models"

type freight int

var Freight freight

func (freight) SaveAll(list []*models.Freight) error {
	if err := models.FreightDefault.Truncate(); err != nil {
		return err
	}
	for _, v := range list {
		if err := v.Save(); err != nil {
			return err
		}
	}
	return nil
}
