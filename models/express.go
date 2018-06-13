package models

type ExpressItem struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}

type ExpressData struct {
	Number string `json:"number"`
	Type   string `json:"type"`
	// 1.在途中 2.正在派件 3.已签收 4.派送失败
	DeliveryStatus string `json:"deliverystatus"`
	// 1.是否签收
	IsSign   string        `json:"issign"`
	ExpName  string        `json:"expName"`
	ExpSite  string        `json:"expSite"`
	ExpPhone string        `json:"expPhone"`
	List     []ExpressItem `json:"list"`
}
