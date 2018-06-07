package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"time"

	"git.jiayougougou.com/snail/transaction/pay/wxpay"
)

func main() {
	bts, err := ioutil.ReadFile("./config_xcx.json")
	if err != nil {
		log.Fatalln(err)
	}
	flag.Parse()
	config := wxpay.PayConfig{}
	if err := json.Unmarshal(bts, &config); err != nil {
		log.Fatalln(err)
	}
	api := wxpay.NewApi(config)
	args := wxpay.NewPayUnifiedOrder()

	args.SetBody("测试商品")
	args.SetOutTradeNo(fmt.Sprintf("TEST%d", time.Now().Unix()))
	args.SetTotalFee(1)
	args.SetOpenId("osh3s0Jd5hsOmlkge0lxvxGXOV90")
	args.SetTradeType("JSAPI")

	result, err := api.UnifiedOrder(args, 5*time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", config)
	fmt.Printf("%+v\n", result)
}
