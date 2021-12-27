package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	config "self/configs"
	dlog2 "self/internal/pkg/dlog"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type DBInit struct {
}

var DbPools []*xorm.Engine

func DbStart() {
	dbConfig := config.Config.Database
	for _, v := range dbConfig.Db.Addr {
		engine, err := xorm.NewEngine("mysql", v)
		if err != nil {
			dlog2.Errorf("connect mysql error, addr: %s, error: %s", v, err.Error())
			continue
		}
		dlog2.Infof("connect mysql ok, addr: %s", v)
		engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
		engine.SetMaxIdleConns(dbConfig.Db.MaxIdle)
		engine.SetMaxOpenConns(dbConfig.Db.MaxIdle)
		//engine.SetTableMapper(names.SameMapper{})   //表名与结构体名称一致
		//engine.SetColumnMapper(names.GonicMapper{}) //结构体驼峰-->表字段名下划线
		if dbConfig.Db.Debug {
			engine.ShowSQL(true)
			//engine.ShowExecTime(true)
			engine.Logger().SetLevel(log.LOG_DEBUG)
		} else {
			engine.Logger().SetLevel(log.LOG_WARNING)
		}
		DbPools = append(DbPools, engine)
	}
	if len(DbPools) == 0 {
		panic("connect mysql all failed")
	}
}

func Db() *xorm.Engine {
	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(len(DbPools))
	return DbPools[randInt]
}

func (this DBInit) DbClose() {
	dlog2.Info("close CloudDb mysql....")
}
