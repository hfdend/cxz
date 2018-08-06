package modules

import (
	"fmt"

	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/utils"
)

type statistics int

var Statistics statistics

type StatisticsData struct {
	TotalVal  float64           `json:"total_val"`
	TrueVal   float64           `json:"true_val"`
	Freight   float64           `json:"freight"`
	PlanVol   float64           `json:"plan_vol"`
	NoPlanVal float64           `json:"no_plan_val"`
	List      []*StatisticsType `json:"list"`
}

type StatisticsType struct {
	Type      string             `json:"type"`
	Number    int                `json:"number"`
	Price     float64            `json:"price"`
	TasteList []*StatisticsTaste `json:"taste_list"`
}

type StatisticsTaste struct {
	Taste  string  `json:"taste"`
	Number int     `json:"number"`
	Price  float64 `json:"price"`
}

func (statistics) Statistics(startTime, endTime int64) (data StatisticsData, err error) {
	var list []*models.Order
	if list, err = models.OrderDefault.GetByTime(startTime, endTime); err != nil {
		return
	}
	storeType := map[string]*StatisticsType{}
	storeTaste := map[string]*StatisticsTaste{}
	for _, v := range list {
		var products []*models.OrderProduct
		if products, err = models.OrderProductDefault.GetByOrderID(v.OrderID); err != nil {
			return
		}
		for _, p := range products {
			if p.IsPlan == models.SureYes {
				data.PlanVol += p.Price
			} else {
				data.NoPlanVal += p.Price
			}
			if _, ok := storeType[p.Type]; !ok {
				t := new(StatisticsType)
				t.Type = p.Type
				data.List = append(data.List, t)
				storeType[p.Type] = t
			}
			ttKey := fmt.Sprintf("%s-%s", p.Type, p.Taste)
			if _, ok := storeTaste[ttKey]; !ok {
				t := new(StatisticsTaste)
				t.Taste = p.Taste
				storeType[p.Type].TasteList = append(storeType[p.Type].TasteList, t)
				storeTaste[ttKey] = t
			}
			storeType[p.Type].Number += p.Number
			storeType[p.Type].Price += utils.Round(p.Price*float64(p.Number), 2)
			storeTaste[ttKey].Number += p.Number
			storeTaste[ttKey].Price += utils.Round(p.Price*float64(p.Number), 2)
		}
		data.TotalVal = v.Price
		data.Freight = v.Freight
		data.TrueVal = v.Price - v.RefundAmount
	}
	return
}
