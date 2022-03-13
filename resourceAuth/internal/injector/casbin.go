package injector

import (
	"github.com/casbin/casbin"
	config2 "self/internal/config"
	"time"

	"github.com/casbin/casbin/persist"
)

// InitCasbin 初始化casbin
func InitCasbin(adapter persist.Adapter) (*casbin.SyncedEnforcer, func(), error) {
	cfg := config2.C.Casbin
	if cfg.Model == "" {
		return new(casbin.SyncedEnforcer), nil, nil
	}

	e := casbin.NewSyncedEnforcer(cfg.Model)

	e.EnableLog(cfg.Debug)

	e.InitWithModelAndAdapter(e.GetModel(), adapter)
	e.EnableEnforce(cfg.Enable)

	cleanFunc := func() {}
	if cfg.AutoLoad {
		e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInternal) * time.Second)
		cleanFunc = func() {
			e.StopAutoLoadPolicy()
		}
	}

	return e, cleanFunc, nil
}
