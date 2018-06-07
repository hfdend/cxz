package wxpay

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// 微信小程序调用支付对象
type PayWxaRequest struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

func (r *PayWxaRequest) SetSign(appId, key string) {
	s := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=%s",
		appId,
		r.NonceStr,
		r.Package,
		r.SignType,
		r.TimeStamp,
		key,
	)
	r.PaySign = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(s))))
}
