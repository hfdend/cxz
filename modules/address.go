package modules

import (
	"strings"

	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
)

type address int

var Address address

func (address) Save(id, userID int, name, phone, code, detailAddress string, isDefault models.Sure) (*models.Address, error) {
	var address *models.Address
	var err error
	if id != 0 {
		if address, err = models.AddressDefault.GetByID(id); err != nil {
			return nil, err
		} else if address == nil {
			return nil, errors.New("地址不存在")
		}
	} else {
		address = new(models.Address)
		address.IsDel = models.SureNo
	}
	switch isDefault {
	case models.SureNo, models.SureYes:
	default:
		return nil, errors.New("是否默认值错误")
	}
	address.IsDefault = isDefault
	address.UserID = userID
	address.Name = name
	address.DistrictCode = code
	districtName, err := models.DistrictDefault.GetNames(code)
	if err != nil {
		return nil, err
	} else if len(districtName) == 0 {
		return nil, errors.New("未知的地址")
	}
	address.DistrictName = strings.Join(districtName, ",")
	address.DetailAddress = detailAddress
	if err := address.Save(); err != nil {
		return nil, err
	}
	if address.IsDefault == models.SureYes {
		if err := models.AddressDefault.UpdateNoDefault(userID, address.ID); err != nil {
			return nil, err
		}
	}
	return address, nil
}
