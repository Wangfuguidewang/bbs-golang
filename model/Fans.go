package model

import (
	"bbs-go/utils/errmsg"
	"errors"
	"gorm.io/gorm"
	"time"
)

// 关注表（Follow）
type Follow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`              // 关注记录ID
	UserID      uint      `gorm:"index" json:"user_id"`              // 发起关注的用户ID
	FollowingID uint      `gorm:"index" json:"following_id"`         // 被关注的用户ID
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`  // 创建时间
	Remarks     string    `gorm:"size:255" json:"remarks,omitempty"` // 关注备注信息，可选字段
}

// 粉丝表（Follower）
type Follower struct {
	ID         uint      `gorm:"primaryKey" json:"id"`             // 粉丝记录ID
	UserID     uint      `gorm:"index" json:"user_id"`             // 用户ID
	FollowerID uint      `gorm:"index" json:"follower_id"`         // 粉丝用户ID
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
}

// 关注follid
func FollwAdd(follid, userid uint) int {
	var foll Follow
	var follower Follower
	code := CheckUserId(int(follid))
	if code == errmsg.ERROR_USER_NOT_EXIST {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	//查询是否关注
	if err := db.Where("following_id = ? AND user_id = ?", follid, userid).First(&foll).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR
		}
	}
	//查询是否是粉丝
	if err := db.Where("follower_id = ? AND user_id = ?", userid, follid).First(&follower).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR
		}
	}
	//添加到1粉丝表和关注表
	//关注
	foll = Follow{
		UserID:      userid,
		FollowingID: follid,
		CreatedAt:   time.Now(),
	}
	//粉丝
	follower = Follower{
		UserID:     follid,
		FollowerID: userid,
		CreatedAt:  time.Now(),
	}
	//关注数量+1
	if err = db.Model(&Profile{}).Where("id = ?", userid).UpdateColumn("following_count", gorm.Expr("following_count + ?", 1)).Error; err != nil {
		return errmsg.ERROR
	}

	//写入到数据库关注表
	err = db.Create(&foll).Error
	if err != nil {
		return errmsg.ERROR
	}
	//粉丝数量+1
	if err = db.Model(&Profile{}).Where("id = ?", follid).UpdateColumn("followers_count", gorm.Expr("followers_count + ?", 1)).Error; err != nil {
		return errmsg.ERROR
	}
	//写入到数据库粉丝表
	err = db.Create(&follower).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 取消关注
func UnFollw(follid, userid uint) int {
	code := CheckUserId(int(follid))
	if code == errmsg.ERROR_USER_NOT_EXIST {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	//删除用户关注
	err = db.Where("following_id = ? AND user_id = ?", follid, userid).Delete(&Follow{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 更新用户关注数
	if err := db.Model(&Profile{}).Where("id = ?", userid).UpdateColumn("following_count", gorm.Expr("following_count - ?", 1)).Error; err != nil {
		return errmsg.ERROR
	}
	//删除用户粉丝
	err = db.Where("follower_id = ? AND user_id = ?", userid, follid).Delete(&Follower{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 更新用户粉丝数
	if err = db.Model(&Profile{}).Where("id = ?", follid).UpdateColumn("followers_count", gorm.Expr("followers_count - ?", 1)).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
