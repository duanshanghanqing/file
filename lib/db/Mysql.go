package db

import (
	"file/lib/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	dbConnect *gorm.DB
)

func createDB(username string, password string, host string, port int64, database string) *gorm.DB {
	//dsn := "gjf:123456@tcp(121.4.213.43:3306)/gjf?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)
	_dbConnect, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true}, // 禁用创建表名复数带 s
		DisableForeignKeyConstraintWhenMigrating: true,                                       // AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
	})
	if err != nil {
		panic("failed to connect database")
	}
	return _dbConnect
}

func init() {
	dbConnect = createDB(config.DB_W_username, config.DB_W_password, config.DB_W_host, config.DB_W_port, config.DB_W_database)
}

func GetDB() *gorm.DB {
	return dbConnect
}
