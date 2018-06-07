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
	bts, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln(err)
	}
	flag.Parse()

	config := wxpay.PayConfig{}
	if err := json.Unmarshal(bts, &config); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(config)

	api := wxpay.NewApi(config)
	args := wxpay.NewPayOrderQuery()
	args.SetOutTradeNo("2017102015281100010426")
	res, err := api.OrderQuery(args, 5*time.Second)
	fmt.Println(err)
	fmt.Printf("%+v\n", res)
}
