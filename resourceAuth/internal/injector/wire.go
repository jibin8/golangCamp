// +build wireinject

package injector

import (
	api2 "self/internal/api"
	bll2 "self/internal/bll/impl/bll"
	"self/internal/model/impl/gorm/model"
	adapter2 "self/internal/module/adapter"
	router2 "self/internal/router"

	"github.com/google/wire"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	// 默认使用gorm存储注入，这里可使用 InitMongoDB & mongoModel.ModelSet 替换为 gorm 存储
	wire.Build(
		InitGormDB,
		model.ModelSet,
		InitAuth,
		InitCasbin,
		InitGinEngine,
		bll2.BllSet,
		api2.APISet,
		router2.RouterSet,
		adapter2.CasbinAdapterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
