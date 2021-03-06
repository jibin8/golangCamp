package internal

import (
	"errors"
	"os"
	"path/filepath"
	config2 "self/internal/config"
	"time"

	"github.com/sirupsen/logrus"
	"self/pkg/logger"
	loggerhook "self/pkg/logger/hook"
	loggergormhook "self/pkg/logger/hook/gorm"
	loggermongohook "self/pkg/logger/hook/mongo"
)

// InitLogger 初始化日志模块
func InitLogger() (func(), error) {
	c := config2.C.Log
	logger.SetLevel(c.Level)
	logger.SetFormatter(c.Format)

	// 设定日志输出
	var file *os.File
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777)

				f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return nil, err
				}
				logger.SetOutput(f)
				file = f
			}
		}
	}

	var hook *loggerhook.Hook
	if c.EnableHook {
		var hookLevels []logrus.Level
		for _, lvl := range c.HookLevels {
			plvl, err := logrus.ParseLevel(lvl)
			if err != nil {
				return nil, err
			}
			hookLevels = append(hookLevels, plvl)
		}

		switch {
		case c.Hook.IsGorm():
			hc := config2.C.LogGormHook

			var dsn string
			switch hc.DBType {
			case "mysql":
				dsn = config2.C.MySQL.DSN()
			case "sqlite3":
				dsn = config2.C.Sqlite3.DSN()
			case "postgres":
				dsn = config2.C.Postgres.DSN()
			default:
				return nil, errors.New("unknown db")
			}

			h := loggerhook.New(loggergormhook.New(&loggergormhook.Config{
				DBType:       hc.DBType,
				DSN:          dsn,
				MaxLifetime:  hc.MaxLifetime,
				MaxOpenConns: hc.MaxOpenConns,
				MaxIdleConns: hc.MaxIdleConns,
				TableName:    hc.Table,
			}),
				loggerhook.SetMaxWorkers(c.HookMaxThread),
				loggerhook.SetMaxQueues(c.HookMaxBuffer),
				loggerhook.SetLevels(hookLevels...),
			)
			logger.AddHook(h)
			hook = h
		case c.Hook.IsMongo():
			h := loggerhook.New(loggermongohook.New(&loggermongohook.Config{
				URI:        config2.C.Mongo.URI,
				Database:   config2.C.Mongo.Database,
				Timeout:    time.Duration(config2.C.Mongo.Timeout) * time.Second,
				Collection: config2.C.LogMongoHook.Collection,
			}),
				loggerhook.SetMaxWorkers(c.HookMaxThread),
				loggerhook.SetMaxQueues(c.HookMaxBuffer),
				loggerhook.SetLevels(hookLevels...),
			)
			logger.AddHook(h)
			hook = h
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}

		if hook != nil {
			hook.Flush()
		}
	}, nil
}
