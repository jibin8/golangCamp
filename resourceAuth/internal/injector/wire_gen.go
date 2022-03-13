//go:generate wire
//+build !wireinject

package injector

import (
	api2 "self/internal/api"
	bll2 "self/internal/bll/impl/bll"
	model2 "self/internal/model/impl/gorm/model"
	adapter2 "self/internal/module/adapter"
	router2 "self/internal/router"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, func(), error) {
	auther, cleanup, err := InitAuth()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := InitGormDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	app := &model2.App{
		DB: db,
	}
	role := &model2.Role{
		DB: db,
	}
	resource := &model2.Resource{
		DB: db,
	}
	roleResource := &model2.RoleResource{
		DB: db,
	}
	account := &model2.Account{
		DB: db,
	}
	accountSystem := &model2.AccountSystem{
		DB: db,
	}
	accountRole := &model2.AccountRole{
		DB: db,
	}
	trans := &model2.Trans{
		DB: db,
	}

	casbinAdapter := &adapter2.CasbinAdapter{
		AccountModel:      	account,
		RoleModel:         	role,
		ResourceModel: 		resource,
		AccountRoleModel:  	accountRole,
		RoleResourceModel:	roleResource,
	}
	syncedEnforcer, cleanup3, err := InitCasbin(casbinAdapter)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}

	// bll
	bllAccount := &bll2.Account{
		Enforcer:			syncedEnforcer,
		TransModel:    		trans,
		AccountModel:     	account,
		AccountSystemModel: accountSystem,
		AccountRoleModel: 	accountRole,
		RoleModel:			role,
	}

	bllLogin := &bll2.Login{
		Enforcer:			syncedEnforcer,
		Auth:            	auther,
		AppModel:			app,
		AccountModel:     	account,
		AccountRoleModel: 	accountRole,
		RoleModel:       	role,
		RoleResourceModel:  roleResource,
		ResourceModel:      resource,
	}

	bllResource := &bll2.Resource{
		TransModel:    		trans,
		ResourceModel:      resource,
	}

	bllApp := &bll2.App{
		AppModel:        	app,
	}

	//bllRole := &bll.Role{
	//	Enforcer:      syncedEnforcer,
	//	TransModel:    trans,
	//	RoleModel:     role,
	//	RoleMenuModel: roleMenu,
	//	UserModel:     user,
	//}
	//
	//bllUser := &bll.User{
	//	Enforcer:      syncedEnforcer,
	//	TransModel:    trans,
	//	UserModel:     user,
	//	UserRoleModel: userRole,
	//	RoleModel:     role,
	//}

	// api
	apiAccount := &api2.Account{
		AccountBll: bllAccount,
	}

	apiLogin := &api2.Login{
		LoginBll: bllLogin,
		AppBll: bllApp,
	}

	//apiMenu := &api.Menu{
	//	MenuBll: bllMenu,
	//}
	//
	//apiRole := &api.Role{
	//	RoleBll: bllRole,
	//}
	//
	//apiUser := &api.Account{
	//	UserBll: bllUser,
	//}

	// router
	routerRouter := &router2.Router{
		Auth:           auther,
		CasbinEnforcer: syncedEnforcer,
		LoginAPI:       apiLogin,
		AccountAPI:     apiAccount,
	}

	engine := InitGinEngine(routerRouter)
	injector := &Injector{
		Engine:         engine,
		Auth:           auther,
		CasbinEnforcer: syncedEnforcer,
		ResourceBll: 	bllResource,
	}
	return injector, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
