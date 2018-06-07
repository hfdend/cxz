package wxpay

import "time"

// 统一下单输入对象
type PayUnifiedOrder interface {
	ToXml() (string, error)
	SetSignType(key SignType)          // 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
	GetSignType() (SignType, bool)     // 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
	SetSign(key string)                // 设置签名，详见签名生成算法
	SetDeviceInfo(device string)       // 设置微信支付分配的终端设备号，商户自定义
	GetDeviceInfo() (string, bool)     // 获取微信支付分配的终端设备号，商户自定义的值
	SetNonceStr(str string)            // 设置随机字符串，不长于32位。推荐随机数生成算法
	GetNonceStr() (string, bool)       // 获取随机字符串，不长于32位。推荐随机数生成算法的值
	SetBody(body string)               // 设置商品或支付单简要描述
	GetBody() (string, bool)           // 获取商品或支付单简要描述的值
	SetDetail(detail string)           // 设置商品名称明细列表
	GetDetail() (string, bool)         // 获取商品名称明细列表的值
	SetAttach(attach string)           // 设置附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	GetAttach() (string, bool)         // 设置附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	SetFeeType(freeType string)        // 设置符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	GetFeeType() (string, bool)        // 获取符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型的值
	SetTotalFee(totalFree int)         // 设置订单总金额，只能为整数，详见支付金额
	GetTotalFee() (int, bool)          // 获取单总金额，只能为整数，详见支付金额
	SetSpbillCreateIp(ip string)       // 设置APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	GetSpbillCreateIp() (string, bool) // 获取APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	SetTimeStart(t time.Time)          // 设置订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	GetTimeStart() (time.Time, bool)   // 获取订单生成时间
	SetTimeExpire(t time.Time)         // 设置订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
	GetTimeExpire() (time.Time, bool)  // 获取订单失效时间
	SetGoodsTag(val string)            // 设置商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	GetGoodsTag() (string, bool)       // 获取商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠的值
	SetNotifyUrl(val string)           // 设置接收微信支付异步通知回调地址
	GetNotifyUrl() (string, bool)      // 获取接收微信支付异步通知回调地址
	SetTradeType(tradeType string)     // 设置取值如下：JSAPI，NATIVE，APP，详细说明见参数规定
	GetTradeType() (string, bool)      // 获取取值如下：JSAPI，NATIVE，APP，详细说明见参数规定的值
	SetProductId(productId string)     // 设置trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
	GetProductId() (string, bool)      // 获取trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。的值
	SetOpenId(openId string)           // 设置trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。下单前需要调用【网页授权获取用户信息】接口获取到用户的Openid。
	GetOpenId() (string, bool)         // 获取trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。下单前需要调用【网页授权获取用户信息】接口获取到用户的Openid。 的值
	SetAppId(appId string)             // 设置微信分配的公众账号ID
	GetAppId() (string, bool)          // 设置微信分配的公众账号ID
	SetMchId(mchId string)             // 设置微信支付分配的商户号
	GetMchId() (string, bool)          // 获取微信支付分配的商户号的值
	SetOutTradeNo(no string)           // 设置商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	GetOutTradeNo() (string, bool)     // 获取商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号的值
	SetAuthCode(authCode string)       // 设置扫码支付授权码，设备读取用户微信中的条码或者二维码信息
	GetAuthCode() (string, bool)       // 获取扫码支付授权码，设备读取用户微信中的条码或者二维码信息的值
}

func NewPayUnifiedOrder() PayUnifiedOrder {
	base := NewPayDataBase()
	return base
}
