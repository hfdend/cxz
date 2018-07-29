package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/hfdend/cxz/cli"
	"github.com/jinzhu/gorm"
)

type Sure int

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
	RedisNil          = redis.Nil
)

const (
	SureNil Sure = 0
	SureYes Sure = 1
	SureNo  Sure = -1
)

// Model 数据库基类
type Model struct {
	db *gorm.DB
	ID int `json:"id" gorm:"primary_key"`
}

// DB 获取数据库client
func (m *Model) DB() *gorm.DB {
	if m.db == nil {
		return cli.DB
	}
	return m.db
}

// SetDB set db
func (m *Model) SetDB(db *gorm.DB) {
	m.db = db
}

// BuildOrderID 生成一个订单ID
func BuildOrderID() (string, error) {
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

func DBInsertIgnore(dbq *gorm.DB, obj interface{}) (int64, error) {
	scope := dbq.NewScope(obj)
	fields := scope.Fields()
	quoted := make([]string, 0, len(fields))
	placeholders := make([]string, 0, len(fields))
	for i := range fields {
		if fields[i].IsPrimaryKey && fields[i].IsBlank {
			continue
		}
		quoted = append(quoted, scope.Quote(fields[i].DBName))
		placeholders = append(placeholders, scope.AddToVars(fields[i].Field.Interface()))
	}

	scope.Raw(fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES (%s)",
		scope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholders, ", ")))

	result, err := scope.SQLDB().Exec(scope.SQL, scope.SQLVars...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func TransactionOver(db *gorm.DB, err error) {
	if db == nil {
		cli.Logger.Error("TransactionOver db is nil")
		return
	}
	if err != nil {
		db.Rollback()
	} else {
		db.Commit()
	}
}

func GetByRedis(c *redis.Client, key string, expiration time.Duration, out interface{}, handler func() ([]byte, error), refreshCaches ...bool) error {
	var refreshCache bool
	if len(refreshCaches) > 0 {
		refreshCache = refreshCaches[0]
	}
	data, err := c.Get(key).Bytes()
	if err == redis.Nil || refreshCache {
		if data, err = handler(); err != nil {
			return err
		} else if len(data) == 0 {
			return nil
		}
		// 保存缓存
		if err := c.Set(key, data, expiration).Err(); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}
