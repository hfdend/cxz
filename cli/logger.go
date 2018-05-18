package cli

import (
	"fmt"
	"log"
	"log/syslog"

	"github.com/Sirupsen/logrus"
	"github.com/hfdend/cxz/conf"
	"github.com/hfdend/cxz/utils"
)

var (
	Logger *logrus.Logger
)

func InitializeLogger() {
	var err error
	Logger, err = getLogger("logrus")
	if err != nil {
		log.Fatalln(err)
	}
}

func getLogger(tag string) (*logrus.Logger, error) {
	loggerConfig := conf.Config.Logger
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000"}
	hk := new(loggerHook)
	network := loggerConfig.Network
	addr := loggerConfig.Addr
	priority := utils.ParseSyslogPriority(loggerConfig.Priority)
	tag = fmt.Sprintf("%s-%s", loggerConfig.PreTag, tag)
	w, err := syslog.Dial(network, addr, priority, tag)
	if err != nil {
		return nil, err
	}
	hk.Writer = w
	logger.Hooks.Add(hk)
	return logger, nil
}

type loggerHook struct {
	Writer *syslog.Writer
}

func (*loggerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (hook *loggerHook) Fire(e *logrus.Entry) error {
	bts, err := (&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:03.000"}).Format(e)
	if err != nil {
		return err
	}
	msg := string(bts)
	switch e.Level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		hook.Writer.Err(msg)
	case logrus.WarnLevel:
		hook.Writer.Warning(msg)
	case logrus.InfoLevel:
		hook.Writer.Info(msg)
	case logrus.DebugLevel:
		hook.Writer.Debug(msg)
	}
	return nil
}
