package middleware

import (
	"colatiger/internal/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitializeDB(conf *viper.Viper) *gorm.DB {
	return initMySqlGorm(conf)
}

// 初始化 mysql gorm.DB
func initMySqlGorm(conf *viper.Viper) *gorm.DB {

	mysqlConfig := mysql.Config{
		DSN:                       conf.GetString("data.mysql"), // DSN data source name
		DefaultStringSize:         191,                          // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                         // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                         // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                         // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                        // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
	}); err != nil {
		log.Printf("mysql connect failed, err:%v", zap.Any("err", err))
		return nil
	} else {
		_, _ = db.DB()
		// 初始化表结构
		initMySqlTables(db)
		return db
	}
}

// 数据库表初始化
func initMySqlTables(db *gorm.DB) {
	log.Print("init table start...")
	err := db.AutoMigrate(
		models.User{},
	)
	if err != nil {
		log.Printf("migrate table failed err:%v", zap.Any("err", err))
		os.Exit(0)
	}
	log.Print("init table end...")
}
