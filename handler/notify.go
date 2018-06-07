package handler

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/hfdend/cxz/conf"
	"github.com/hfdend/cxz/modules"
	"github.com/hfdend/cxz/payment/wxpay"
)

var logger = logrus.New()

func WXAPaymentNotify(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		wxNotifyReturn(c, fmt.Errorf("not found attach"))
		return
	}

	logger.Infof("wexin notify: receive msg: %s", string(body))
	var values = wxpay.DataValues{}
	if err := xml.Unmarshal(body, &values); err != nil {
		wxNotifyReturn(c, err)
		return
	}
	data := wxpay.NewPayDataBase()
	data.Values = values
	var wxc wxpay.PayConfig
	wxc.AppId = conf.Config.WXPay.AppId
	wxc.MchId = conf.Config.WXPay.MchId
	wxc.Key = conf.Config.WXPay.Key
	wxc.NotifyUrl = conf.Config.WXPay.NotifyUrl
	api := wxpay.NewApi(wxc)
	api.Logger = logger
	// 验证返回签名
	sign := data.MakeSign(api.Config.Key)
	if sign != data.GetSign() {
		wxNotifyReturn(c, fmt.Errorf("签名错误"))
		return
	}
	// 判断支付结果
	if data.GetReturnCode() != "SUCCESS" || data.GetResultCode() != "SUCCESS" {
		// 如果没有支付成功则不做任何操作
		wxNotifyReturn(c, "OK")
		return
	}
	// 获取订单号
	orderID, ok := data.GetOutTradeNo()
	if !ok {
		wxNotifyReturn(c, fmt.Errorf("OutTradeNo is empty"))
		return
	}
	transactionID, ok := data.GetTransactionId()
	if !ok {
		wxNotifyReturn(c, fmt.Errorf("TransactionId is empty"))
		return
	}
	// 更改订单到支付成功
	if err := modules.Order.PaymentSuccess(orderID, transactionID); err != nil {
		wxNotifyReturn(c, err)
		return
	}
	wxNotifyReturn(c, "success")
}

func wxNotifyReturn(c *gin.Context, data interface{}) {
	t := template.Must(template.New("name").Parse(`<xml><return_code><![CDATA[{{.Code}}]]></return_code><return_msg><![CDATA[{{.Msg}}]]></return_msg></xml>`))
	var mp struct {
		Code string
		Msg  string
	}
	defer func() {
		logger.Info("weixin pay return")
	}()
	if err, ok := data.(error); ok {
		// 记录错误
		logger.WithError(err).Error("weixin pay notify error")
		mp.Code = "FAIL"
		mp.Msg = err.Error()
	} else {
		mp.Code = "SUCCESS"
		mp.Msg = "OK"
	}
	t.Execute(c.Writer, mp)
	return
}
