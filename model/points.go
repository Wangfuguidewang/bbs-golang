package model

import (
	"bbs-go/utils/errmsg"
	"gorm.io/gorm"
)

// 积分+1
func PointsAdd(userid uint) int {
	var user User
	db.Select("points").Find(&user)
	err := db.Model(&User{}).Where("id = ?", userid).UpdateColumn("points", gorm.Expr("points + ?", 1)).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 积分-1
func PointsDe(userid uint) int {
	var user User
	db.Select("points").Where("id = ?", userid).First(&user)
	if user.Points <= 0 {
		return errmsg.ERROR
	}
	err := db.Model(&User{}).Where("id = ?", userid).UpdateColumn("points", gorm.Expr("points - ?", 1)).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
