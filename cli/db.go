package cli

import (
	"database/sql/driver"
	"fmt"
	"log"
	"log/syslog"
	"reflect"
	"regexp"
	"time"
	"unicode"

	"github.com/go-sql-driver/mysql"
	"github.com/hfdend/cxz/conf"
	"github.com/hfdend/cxz/utils"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitMysql() {
	mysqlConfig := conf.Config.Mysql
	sourceName := (&mysql.Config{
		AllowNativePasswords: true,
		User:                 mysqlConfig.User,
		Passwd:               mysqlConfig.Password,
		Addr:                 fmt.Sprintf("%s:%d", mysqlConfig.Host, mysqlConfig.Port),
		Net:                  "tcp",
		Params: map[string]string{
			"charset": "utf8",
			"loc":     "Local",
		},
		DBName:  mysqlConfig.DB,
		Timeout: mysqlConfig.Timeout,
	}).FormatDSN()
	var err error
	if DB, err = gorm.Open("mysql", sourceName); err != nil {
		log.Fatalln(err, sourceName)
	}
	DB.DB().SetConnMaxLifetime(time.Duration(mysqlConfig.ConnMaxLifetime) * time.Second)
	DB.DB().SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	DB.DB().SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	// 设置logger
	DB.LogMode(true)
	network := conf.Config.Logger.Network
	addr := conf.Config.Logger.Addr
	priority := utils.ParseSyslogPriority(conf.Config.Logger.Priority)
	tag := fmt.Sprintf("%s-sql", conf.Config.Logger.PreTag)
	w, err := syslog.Dial(network, addr, priority, tag)
	if err != nil {
		log.Fatalln(err)
	}
	logger := &DBLogger{Write: w}
	DB.SetLogger(logger)
}

var (
	sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)
)

// Logger default logger
type DBLogger struct {
	Write *syslog.Writer
}

type loggerData struct {
	Source  string
	Time    time.Time
	Level   string
	UseTime string
	Sql     string
	Message string
}

type sqlLogger struct {
	w *syslog.Writer
}

func (ld loggerData) ToString() string {
	return fmt.Sprintf("%+v", ld)
}

// Print format & print log
func (logger DBLogger) Print(values ...interface{}) {
	if len(values) > 1 {
		var data loggerData
		level := values[0]
		data.Source = fmt.Sprintf("%v", values[1])
		data.Time = time.Now()
		data.Level = fmt.Sprintf("%v", level)
		if level == "sql" {
			data.UseTime = values[2].(time.Duration).String()
			var sql string
			var formattedValues []string
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}

			var formattedValuesLength = len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
			data.Sql = sql
			logger.Write.Info(data.ToString())
		} else {
			data.Message = fmt.Sprint(values[2:]...)
			logger.Write.Err(data.ToString())
		}
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
