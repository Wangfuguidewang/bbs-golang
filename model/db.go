package model

import (
	"bbs-go/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	db  = new(gorm.DB)
	err error
)

// 连接数据库
func Initdb() {
	db, err = gorm.Open(mysql.Open(
		fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
			utils.DbUser,
			utils.DbPassWord,
			utils.DbHost,
			utils.DbPort,
			utils.DbName),
	))
	if err != nil {
		fmt.Println("连接失败", err)
		fmt.Println(utils.DbUser,
			utils.DbPassWord,
			utils.DbHost,
			utils.DbPort,
			utils.DbName)
		return
	}
	fmt.Println("连接成功")
	db.AutoMigrate(&User{}, &Profile{}, &Topic{}, &Comment{}, &Like{}, &Favorite{})
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, _ := db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

}
