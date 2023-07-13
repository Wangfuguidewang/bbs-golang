package model

import (
	"bbs-go/utils/errmsg"
	"gorm.io/gorm"
	"time"
)

type Topic struct { //话题
	gorm.Model
	UserID  uint   `gorm:"comment:用户ID" json:"user_id"` // 用户ID
	Title   string `gorm:"comment:标题" json:"title"`     // 标题
	Content string `gorm:"comment:内容" json:"content"`   // 内容

}

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`  // 点赞ID，主键
	TopicID   uint      `gorm:"index" json:"topic_id"` // 所属话题ID，用于与话题表进行关联
	UserID    uint      `gorm:"index" json:"user_id"`  // 点赞用户ID，用于与用户表进行关联
	CreatedAt time.Time `json:"created_at"`            // 创建时间
}

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`  // 收藏ID，主键
	TopicID   uint      `gorm:"index" json:"topic_id"` // 所属话题ID，用于与话题表进行关联
	UserID    uint      `gorm:"index" json:"user_id"`  // 收藏用户ID，用于与用户表进行关联
	CreatedAt time.Time `json:"created_at"`            // 创建时间
}

// 新增话题
func CreateTopic(data *Topic) int {
	var comment Comment
	var like Like
	var favorite Favorite
	code := CheckUserId(int(data.UserID))
	if code == 1003 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	comment.TopicID = data.ID
	err = db.Create(&comment).Error
	if err != nil {
		return errmsg.ERROR
	}
	like.TopicID = data.ID
	err = db.Create(&like).Error
	if err != nil {
		return errmsg.ERROR
	}
	favorite.TopicID = data.ID
	err = db.Create(&favorite).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// EditTopic 编辑话题
func EditTopic(id int, userid uint, data *Topic) int {
	var art Topic
	art, _ = GetTopInfo(id)
	if art.UserID != userid {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["user_id"] = userid
	maps["content"] = data.Content

	err = db.Model(&art).Where("id = ? ", id).Updates(&maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//GetCaTopic查询userid下所有话题

func GetCaTopic(id int, pageSize int, pageNum int) ([]Topic, int, int64) {
	var cateTopic []Topic
	var total int64
	err = db.Limit(pageSize).Offset((pageNum-1)*pageSize).Where(
		"user_id =?", id).Find(&cateTopic).Error
	db.Model(&cateTopic).Where("user_id =?", id).Count(&total)
	if err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}
	return cateTopic, errmsg.SUCCSE, total
}

// GetTopInfo 查询单个话题
func GetTopInfo(id int) (Topic, int) {
	var top Topic
	err = db.Where("id = ?", id).First(&top).Error
	if err != nil {
		return top, errmsg.ERROR_ART_NOT_EXIST
	}
	return top, errmsg.SUCCSE
}

// DeleteTop 删除话题
func DeleteTop(id int, userid uint) int {
	var art Topic
	art, code := GetTopInfo(id)
	if art.UserID != userid {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	if code == errmsg.ERROR_ART_NOT_EXIST {
		return errmsg.ERROR_ART_NOT_EXIST
	}
	err = db.Where("id = ? ", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}

	return errmsg.SUCCSE
}
