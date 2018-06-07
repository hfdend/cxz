package wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type CancelResults struct {
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
	Recall     string `xml:"recall"`       // 是否需要继续调用撤销，Y-需要，N-不需要
}

func (res CancelResults) JudgmentResultCode() error {
	if res.ResultCode != "SUCCESS" {
		return errors.New(res.ReturnMsg)
	}
	if res.ResultCode != "SUCCESS" {
		return errors.New(fmt.Sprintf("%s: %s", res.ErrCode, res.ErrCodeDes))
	}
	return nil
}
