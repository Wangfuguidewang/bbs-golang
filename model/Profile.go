package model

import "bbs-go/utils/errmsg"

type Profile struct {
	ID             uint   `gorm:"primarykey"`
	Username       string `gorm:"not null;comment:用户名"json:"username"`
	Nickname       string `gorm:"comment:昵称" json:"nickname"`       // 昵称
	Email          string `gorm:"comment:邮箱"json:"email"`           //邮箱
	Bio            string `gorm:"comment:用户简介" json:"bio"`          // 用户简介
	Avatar         string `gorm:"comment:个人头像" json:"avatar"`       // 个人头像
	FollowersCount uint   `gorm:"default:0" json:"followers_count"` // 粉丝数量
	FollowingCount uint   `gorm:"default:0" json:"following_count"` // 关注数量
}

// Getprofile 获取个人信息
func Getprofile(id int) (Profile, int) {
	var profile Profile
	err = db.Where("ID = ?", id).First(&profile).Error
	if err != nil {
		return profile, errmsg.ERROR
	}
	return profile, errmsg.SUCCESS
}

//updateProfile 更新个人信息

func UpdateProfile(id int, data *Profile) int {
	var profile Profile
	var user User
	err = db.Model(&profile).Where("ID = ?", id).Updates(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	user.UserName = data.Username
	user.Email = data.Email
	err = db.Model(&user).Where("ID = ?", id).Updates(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
