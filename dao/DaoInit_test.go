package dao

import (
	"file/lib/db"
	"file/model"
	"testing"
)

// go test -v -run TestDaoInit_CreateFileTable DaoInit_test.go
func TestDaoInit_CreateFileTable(t *testing.T) {
	_db := db.GetDB()

	// 创建file表
	err := _db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.File{})
	if err != nil {
		t.Logf(`创建file表失败`)
		return
	}
	t.Logf(`创建file表成功`)
}
