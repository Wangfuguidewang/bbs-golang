package model

import (
	"bbs-go/utils/errmsg"
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`     // 评论ID，主键
	TopicID   uint      `gorm:"index" json:"topic_id"`    // 所属话题ID，用于与话题表进行关联
	UserID    uint      `gorm:"index" json:"user_id"`     // 评论用户ID，用于与用户表进行关联
	Content   string    `gorm:"type:text" json:"content"` // 评论内容
	CreatedAt time.Time `json:"created_at"`               // 创建时间
	UpdatedAt time.Time `json:"updated_at"`               // 更新时间
}

// AddComment 新增评论
func AddComment(data *Comment) int {
	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE

}
