package conf

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Mysql struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DB              string        `yaml:"db"`
	Timeout         time.Duration `yaml:"timeout"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
}

type Redis struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	DB          int           `yaml:"db"`
	PoolSize    int           `yaml:"pool_size"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

var Config struct {
	Main struct {
		Mode string `yaml:"mode"`
		Addr string `yaml:"addr"`
	} `yaml:"main"`
	Mysql  Mysql `yaml:"mysql"`
	Redis  Redis `yaml:"redis"`
	Logger struct {
		Network  string `yaml:"network"`
		Addr     string `yaml:"addr"`
		Priority string `yaml:"priority"`
		PreTag   string `yaml:"pre_tag"`
	} `yaml:"logger"`
	Aliyun struct {
		OSS struct {
			Domain          string `yaml:"domain"`
			Endpoint        string `yaml:"endpoint"`
			AccessKeyID     string `yaml:"access_key_id"`
			AccessKeySecret string `yaml:"access_key_secret"`
			Bucket          string `yaml:"bucket"`
		} `yaml:"oss"`
		SMS struct {
			Test            bool   `yaml:"test"`
			AccessKeyID     string `yaml:"access_key_id"`
			AccessKeySecret string `yaml:"access_key_secret"`
		} `yaml:"sms"`
		Express struct {
			AppCode string `yaml:"app_code"`
		} `yaml:"express"`
	} `yaml:"aliyun"`
	MiniProgram struct {
		AppID  string `yaml:"appid"`
		Secret string `yaml:"secret"`
	} `yaml:"miniprogram"`
	WXPay struct {
		AppId      string  `yaml:"appid"`
		MchId      string  `yaml:"mchid"`
		Key        string  `yaml:"key"`
		NotifyUrl  string  `yaml:"notify_url"`
		TestAmount float64 `yaml:"test_amount"`
	} `yaml:"wxpay"`
}

type FlagParseFn = func(fs *flag.FlagSet)

func Init(fns ...FlagParseFn) {
	var configFile string
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	flag.StringVar(&configFile, "f", "", "config file")
	for _, f := range fns {
		f(flag.CommandLine)
	}
	flag.Parse()
	var (
		data []byte
		err  error
	)
	if configFile == "" {
		if _, err := os.Stat("./config.yml"); err == nil {
			configFile = "./config.yml"
		}
	}
	if configFile == "" {
		if _, err := os.Stat("../config.yml"); err == nil {
			configFile = "../config.yml"
		}
	}
	if configFile == "" {
		if _, err := os.Stat("../../config.yml"); err == nil {
			configFile = "../../config.yml"
		}
	}
	if configFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if data, err = ioutil.ReadFile(configFile); err != nil {
		log.Fatalln(err)
	}
	if err = yaml.Unmarshal(data, &Config); err != nil {
		log.Fatalln(err)
	}
}
