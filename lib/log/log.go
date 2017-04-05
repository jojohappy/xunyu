package log

import (
	"fmt"
	golog "log"
	"os"
	"path"

	"github.com/xunyu/config"
)

const logFlag = golog.Ldate | golog.Ltime | golog.Lmicroseconds | golog.Lshortfile

type Priority int

const (
	LOG_CRIT Priority = iota
	LOG_ERR
	LOG_WARNING
	LOG_INFO
	LOG_DEBUG
)

type LogConfig struct {
	File  string `config:"file"`
	Level int    `config:"level"`
}

type Logger struct {
	file   *os.File
	logger *golog.Logger
	level  Priority
	config LogConfig
}

var (
	defaultLogConfig = LogConfig{
		File:  "logs/xunyu.log",
		Level: 3,
	}
	_log = Logger{}
)

func InitLog(cfg *config.Config) error {
	logConfig := defaultLogConfig
	if err := cfg.Assemble(&logConfig); nil != err {
		return err
	}

	logDir := path.Dir(logConfig.File)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0750)
		if err != nil {
			return err
		}
	}

	logFile, err := os.OpenFile(logConfig.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		golog.Fatalf("open log file error: %s\n", err)
	}
	_log.config = logConfig
	_log.file = logFile
	_log.logger = golog.New(_log.file, "", logFlag)
	_log.level, err = getLogLevel(logConfig.Level)
	if nil != err {
		return err
	}
	return nil
}

func Stop() {
	_log.file.Close()
}

func getLogLevel(l int) (Priority, error) {
	levels := map[int]Priority{
		0: LOG_CRIT,
		1: LOG_ERR,
		2: LOG_WARNING,
		3: LOG_INFO,
		4: LOG_DEBUG,
	}

	level, ok := levels[l]
	if !ok {
		return 0, fmt.Errorf("unknown log level: %v", l)
	}
	return level, nil
}

func marked(level Priority, prefix string, format string, v ...interface{}) {
	if _log.level >= level {
		s := fmt.Sprintf(prefix+format+"\n", v...)
		_log.logger.Output(3, s)
	}
}

func Debug(format string, v ...interface{}) {
	marked(LOG_DEBUG, "DEBUG ", format, v...)
}

func Info(format string, v ...interface{}) {
	marked(LOG_INFO, "INFO ", format, v...)
}

func Warn(format string, v ...interface{}) {
	marked(LOG_WARNING, "WARNING ", format, v...)
}

func Err(format string, v ...interface{}) {
	marked(LOG_ERR, "ERROR ", format, v...)
}

func Critical(format string, v ...interface{}) {
	marked(LOG_CRIT, "CRITICAL ", format, v...)
}
