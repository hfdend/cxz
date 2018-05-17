package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hfdend/cxr/cli"
	"github.com/jinzhu/gorm"
)

type model struct {
	db *gorm.DB
}

func (m *model) DB() *gorm.DB {
	if m.db == nil {
		return cli.DB
	} else {
		return m.db
	}
}

func (m *model) SetDB(db *gorm.DB) {
	m.db = db
}

func BuildOrderId() (string, error) {
	// 当前时间格式
	now := time.Now().Format("20060102150405")
	//  4位随机数
	rd := fmt.Sprintf("%0.4d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
	// redis自增保证唯一
	key := fmt.Sprintf("%s%s%s", "transaction_id_inc", now, rd)
	id, err := cli.Redis.Incr(key).Result()
	if err != nil {
		return "", err
	}
	if err := cli.Redis.Expire(key, 2*time.Second).Err(); err != nil {
		return "", err
	}
	no := fmt.Sprintf("%s%0.4d%s", now, id, rd)
	return no, nil
}
