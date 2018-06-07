package wxpay

import (
	"encoding/xml"
	"errors"
)

// 接口调用结果类
type PayResults struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"` // 此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
	ReturnMsg  string   `xml:"return_msg"`  // 返回信息，如非空，为错误原因 签名失败 参数格式校验错误
	// 当return_code为SUCCESS的时候，还会包括以下字段：
	AppId      string `xml:"appid"`        // 调用接口提交的公众账号ID
	MchId      string `xml:"mch_id"`       // 调用接口提交的商户号
	DeviceInfo string `xml:"device_info"`  // 调用接口提交的终端设备号
	NonceStr   string `xml:"nonce_str"`    // 微信返回的随机字符串
	Sign       string `xml:"sign"`         // 微信返回的签名
	ResultCode string `xml:"result_code"`  // SUCCESS/FAIL
	ErrCode    string `xml:"err_code"`     // 详细参见错误列表
	ErrCodeDes string `xml:"err_code_des"` // 错误返回的信息描述
	// 当return_code 和result_code都为SUCCESS的时，还会包括以下字段：
	// ---扫码支付字段
	PrepayId string `xml:"prepay_id"` // 微信生成的预支付会话标识，用于后续接口调用中使用，该值有效期为2小时
	CodeUrl  string `xml:"code_url"`  // trade_type为NATIVE时有返回，用于生成二维码，展示给用户进行扫码支付

	OpenId             string `xml:"openid"`               // 用户在商户appid 下的唯一标识
	IsSubscribe        string `xml:"is_subscribe"`         // 用户是否关注公众账号，仅在公众账号类型支付有效，取值范围：Y或N;Y-关注;N-未关注
	TradeType          string `xml:"trade_type"`           // 支付类型为MICROPAY(即扫码支付)
	BankType           string `xml:"bank_type"`            // 银行类型，采用字符串类型的银行标识
	FeeType            string `xml:"fee_type"`             // 符合ISO 4217标准的三位字母代码，默认人民币：CNY
	TotalFee           string `xml:"total_fee"`            // 订单总金额，单位为分，只能为整数，详见支付金额
	SettlementTotalFee string `xml:"settlement_total_fee"` // 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	CouponFee          string `xml:"coupon_fee"`           // “代金券”金额<=订单金额，订单金额-“代金券”金额=现金支付金额，详见支付金额
	CashFeeType        string `xml:"cash_fee_type"`        // 符合ISO 4217标准的三位字母代码，默认人民币：CNY
	CashFee            string `xml:"cash_fee"`             // 订单现金支付金额
	TransactionId      string `xml:"transaction_id"`       // 微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	Attach             string `xml:"attach"`               // 商家数据包，原样返回
	TimeEnd            string `xml:"time_end"`             // 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010
	PromotionDetail    string `xml:"promotion_detail"`     // 新增返回，单品优惠功能字段，需要接入请见
}

func (result PayResults) JudgmentResultCode() (string, error) {
	switch result.ErrCode {
	default:
		return "未知错误码", errors.New("unknown result code")
	case "SYSTEMERROR":
		return "系统超时", nil // 请立即调用被扫订单结果查询API，查询当前订单状态，并根据订单的状态决定下一步的操作。
	case "PARAM_ERROR":
		return "参数错误", errors.New("参数错误")
	case "ORDERPAID":
		return "订单已支付", errors.New("订单已支付")
	case "NOAUTH":
		return "商户无权限", errors.New("商户无权限")
	case "AUTHCODEEXPIRE":
		return "二维码已过期，请用户在微信上刷新后再试", nil
	case "NOTENOUGH":
		return "余额不足", errors.New("余额不足")
	case "NOTSUPORTCARD":
		return "不支持卡类型", errors.New("不支持卡类型")
	case "ORDERCLOSED":
		return "订单已关闭", errors.New("订单已关闭")
	case "ORDERREVERSED":
		return "订单已撤销", errors.New("订单已撤销")
	case "BANKERROR":
		return "银行系统异常", nil
	case "USERPAYING":
		return "用户支付中，需要输入密码", nil
	case "AUTH_CODE_ERROR":
		return "授权码参数错误", errors.New("授权码参数错误")
	case "AUTH_CODE_INVALID":
		return "授权码检验错误", errors.New("授权码检验错误")
	case "XML_FORMAT_ERROR":
		return "XML格式错误", errors.New("XML格式错误")
	case "REQUIRE_POST_METHOD":
		return "请使用post方法", errors.New("请使用post方法")
	case "SIGNERROR":
		return "签名错误", errors.New("签名错误")
	case "LACK_PARAMS":
		return "缺少参数", errors.New("缺少参数")
	case "NOT_UTF8":
		return "编码格式错误", errors.New("编码格式错误")
	case "BUYER_MISMATCH":
		return "支付帐号错误", errors.New("支付帐号错误")
	case "APPID_NOT_EXIST":
		return "APPID不存在", errors.New("APPID不存在")
	case "MCHID_NOT_EXIST":
		return "MCHID不存在", errors.New("MCHID不存在")
	case "OUT_TRADE_NO_USED":
		return "商户订单号重复", errors.New("商户订单号重复")
	case "APPID_MCHID_NOT_MATCH":
		return "appid和mch_id不匹配", errors.New("appid和mch_id不匹配")
	case "INVALID_REQUEST":
		return "无效请求", errors.New("无效请求")
	case "TRADE_ERROR":
		return "交易错误", errors.New("交易错误")
	}
}
