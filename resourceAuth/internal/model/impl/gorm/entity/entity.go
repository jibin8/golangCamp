package entity

import (
	"context"
	"fmt"
	config2 "self/internal/config"
	icontext2 "self/internal/icontext"
	"time"

	"github.com/jinzhu/gorm"
)

// Model base model
type Model struct {

}

// CommonModel base model
type CommonModel struct {
	ID        	int   		`gorm:"column:id;primary_key;size:11;"`
	CreateTime 	time.Time  	`gorm:"column:createTime;"`
	UpdateTime 	time.Time  	`gorm:"column:updateTime;"`
}

// TableName table name
func (Model) TableName(name string) string {
	return fmt.Sprintf("%s%s", config2.C.Gorm.TablePrefix, name)
}

// GetDB ...
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := icontext2.FromTrans(ctx)
	if ok && !icontext2.FromNoTrans(ctx) {
		db, ok := trans.(*gorm.DB)
		if ok {
			if icontext2.FromTransLock(ctx) {
				if dbType := config2.C.Gorm.DBType; dbType == "mysql" ||
					dbType == "postgres" {
					db = db.Set("gorm:query_option", "FOR UPDATE")
				}
			}
			return db
		}
	}
	return defDB
}

// GetDBWithModel ...
func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}
