package wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type QueryResults struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"` // SUCCESS/FAIL 此字段是通信标识，非交易标识，交易是否成功需要查看trade_state来判断
	ReturnMsg  string   `xml:"return_msg"`  // 返回信息，如非空，为错误原因 签名失败 参数格式校验错误
	// 以下字段在return_code为SUCCESS的时候有返回
	AppId      string `xml:"appid"`        // 微信分配的公众账号ID
	MchId      string `xml:"mch_id"`       // 微信支付分配的商户号
	NonceStr   string `xml:"nonce_str"`    // 随机字符串，不长于32位
	Sign       string `xml:"sign"`         // 签名，详见
	ResultCode string `xml:"result_code"`  // SUCCESS/FAIL
	ErrCode    string `xml:"err_code"`     // 错误码
	ErrCodeDes string `xml:"err_code_des"` // 结果信息描述
	// 以下字段在return_code 、result_code、trade_state都为SUCCESS时有返回 ，如trade_state不为 SUCCESS，则只返回out_trade_no（必传）和attach（选传）。
	DeviceInfo         string `xml:"device_info"`          // 微信支付分配的终端设备号
	OpenId             string `xml:"openid"`               // 用户在商户appid下的唯一标识
	IsSubscribe        string `xml:"is_subscribe"`         // 用户是否关注公众账号，Y-关注，N-未关注，仅在公众账号类型支付有效
	TradeType          string `xml:"trade_type"`           // 调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，MICROPAY，详细说明见参数规定
	TradeState         string `xml:"trade_state"`          // SUCCESS—支付成功 REFUND—转入退款 NOTPAY—未支付 CLOSED—已关闭 REVOKED—已撤销（刷卡支付） USERPAYING--用户支付中 PAYERROR--支付失败(其他原因，如银行返回失败) 支付状态机请见下单API页面
	BankType           string `xml:"bank_type"`            // 银行类型，采用字符串类型的银行标识
	TotalFee           int64  `xml:"total_fee"`            // 订单总金额，单位为分
	SettlementTotalFee int64  `xml:"settlement_total_fee"` // 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	FeeType            string `xml:"fee_type"`             // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY
	CashFee            int64  `xml:"cash_fee"`             // 现金支付金额订单现金支付金额，详见
	CashFeeType        string `xml:"cash_fee_type"`        // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY
	TransactionId      string `xml:"transaction_id"`       // 微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	Attach             string `xml:"attach"`               // 附加数据，原样返回
	TimeEnd            string `xml:"time_end"`             // 订单支付时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见
	TradeStateDesc     string `xml:"trade_state_desc"`     // 对当前查询订单状态的描述和下一步操作的指引
}

func (res QueryResults) JudgmentResultCode() error {
	if res.ResultCode != "SUCCESS" {
		return errors.New(res.ReturnMsg)
	}
	if res.ResultCode != "SUCCESS" {
		return errors.New(fmt.Sprintf("%s: %s", res.ErrCode, res.ErrCodeDes))
	}
	return nil
}
