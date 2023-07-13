package model

import (
	"bbs-go/utils/errmsg"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Topic struct { //话题
	gorm.Model
	UserID     uint   `gorm:"comment:用户ID" json:"user_id"`     // 用户ID
	Title      string `gorm:"comment:标题" json:"title"`         // 标题
	Content    string `gorm:"comment:内容" json:"content"`       // 内容
	Likes      []Like `gorm:"foreignKey:TopicID" json:"likes"` // 文章点赞记录
	LikesCount uint   `gorm:"default:0" json:"likes_count"`    // 文章点赞数
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
	/*var favorite Favorite*/
	code := CheckUserId(int(data.UserID))
	if code == 1003 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	/*comment.TopicID = data.ID
	err = db.Create(&comment).Error
	if err != nil {
		return errmsg.ERROR
	}
	like.TopicID = data.ID
	err = db.Create(&like).Error
	if err != nil {
		return errmsg.ERROR
	}*/
	//favorite.TopicID = data.ID
	//err = db.Create(&favorite).Error
	//if err != nil {
	//	return errmsg.ERROR
	//}
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

// 点赞
func LikeTopic(topicID, userID uint) int {
	// 检查是否已经点赞过
	var like Like
	_, err := GetTopInfo(int(topicID))
	if err != errmsg.SUCCSE {
		return errmsg.ERROR_ART_NOT_EXIST
	}
	if err := db.Where("topic_id = ? AND user_id = ?", topicID, userID).First(&like).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR
		}

		// 创建点赞记录
		like = Like{
			TopicID:   topicID,
			UserID:    userID,
			CreatedAt: time.Now(),
		}

		// 更新文章点赞数
		if err := db.Model(&Topic{}).Where("id = ?", topicID).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1)).Error; err != nil {
			return errmsg.ERROR
		}

		// 保存点赞记录到数据库
		if err := db.Create(&like).Error; err != nil {
			return errmsg.ERROR
		}
	}

	return errmsg.SUCCSE
}

// 取消点赞
func Unliketop(topicId, userId uint) int {
	_, code := GetTopInfo(int(topicId))
	if code != errmsg.SUCCSE {
		return errmsg.ERROR_ART_NOT_EXIST
	}
	//删除点赞数
	err = db.Where("topic_id = ? AND user_id = ?", topicId, userId).Delete(&Like{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	// 更新文章点赞数
	if err := db.Model(&Topic{}).Where("id = ?", topicId).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1)).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 添加收藏
func FavoriTopADD(topid, userid uint) int {
	var favo Favorite
	_, err := GetTopInfo(int(topid))
	if err != errmsg.SUCCSE {
		return errmsg.ERROR_ART_NOT_EXIST
	}
	if err := db.Where("topic_id = ? AND user_id = ?", topid, userid).First(&favo).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.ERROR
		}

		//
		favo = Favorite{
			TopicID:   topid,
			UserID:    userid,
			CreatedAt: time.Now(),
		}
		if err := db.Create(&favo).Error; err != nil {
			return errmsg.ERROR
		}
	}
	return errmsg.SUCCSE
}

// 删除收藏
func DeFavoTop(topid, userid uint) int {
	_, code := GetTopInfo(int(topid))
	if code != errmsg.SUCCSE {
		return errmsg.ERROR_ART_NOT_EXIST
	}
	//删除收藏数
	err = db.Where("topic_id = ? AND user_id = ?", topid, userid).Delete(&Favorite{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
