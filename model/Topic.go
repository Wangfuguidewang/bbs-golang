package model

import (
	"bbs-go/utils/errmsg"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sort"
	"time"
)

type Topic struct { //话题
	gorm.Model
	UserID     uint   `gorm:"comment:用户ID" json:"user_id"`     // 用户ID
	Title      string `gorm:"comment:标题" json:"title"`         // 标题
	Content    string `gorm:"comment:内容" json:"content"`       // 内容
	Likes      []Like `gorm:"foreignKey:TopicID" json:"likes"` // 文章点赞记录
	LikesCount int    `gorm:"default:0" json:"likes_count"`    // 文章点赞数
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
	//积分+1
	PointsAdd(data.UserID)

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
	return errmsg.SUCCESS
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
	return errmsg.SUCCESS
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
	return cateTopic, errmsg.SUCCESS, total
}

// // 查找所有话题
//
//	func GetTopic() ([]Topic, int) {
//		// 初始化 Redis 客户端
//		client := redis.NewClient(&redis.Options{
//			Addr: "localhost:6379", // Redis 服务器地址
//		})
//		var topics []Topic
//		err = db.Find(&topics).Error
//		if err != nil {
//			return nil, errmsg.ERROR
//		}
//
//		for i := range topics {
//			fmt.Println(topics[i].LikesCount)
//			if topics[i].LikesCount < topics[i+1].LikesCount {
//				topics[i], topics[i+1] = topics[i+1], topics[i]
//
//			}
//		}
//		// 将话题列表存储到 Redis 的有序集合中
//		ctx := context.Background()
//		var Topicc []Topic
//		for _, topic := range topics {
//			_, err := json.Marshal(topic)
//			if err != nil {
//				return nil, errmsg.ERROR
//			}
//			for _, topic := range topics {
//				_, err := json.Marshal(topic)
//				if err != nil {
//					panic(err)
//				}
//			}
//			// 设置 Redis 有序集合的过期时间为五秒
//			client.Expire(ctx, "topics", 5*time.Second)
//
//			// 从 Redis 获取排序后的话题列表并返回给客户端
//			sortedTopics, err := client.ZRevRangeByScore(ctx, "topics", &redis.ZRangeBy{
//				Min: "-inf",
//				Max: "+inf",
//			}).Result()
//			if err != nil {
//				panic(err)
//			}
//
//			for i, topicJSON := range sortedTopics {
//				var topic Topic
//				err = json.Unmarshal([]byte(topicJSON), &topic)
//				if err != nil {
//					return nil, errmsg.ERROR
//				}
//				// 输出话题信息，这里可以根据需要自行处理
//				Topicc[i] = topic
//				fmt.Println(topic.LikesCount)
//			}
//		}
//		return Topicc, errmsg.SUCCESS
//	}
//
// 模拟从 MySQL 查询数据的函数
func getDataFromMySQL(orderBy string, asc bool) ([]Topic, int) {
	var dataList []Topic
	err := db.Find(&dataList).Error
	if err != nil {
		return nil, errmsg.ERROR
	}

	//按照字段排序（这里以 ID 字段为例）
	if orderBy == "likes_count" {
		if asc {
			// 升序排序
			sort.Slice(dataList, func(i, j int) bool {
				return dataList[i].LikesCount < dataList[j].LikesCount
			})
		} else {
			// 降序排序
			sort.Slice(dataList, func(i, j int) bool {
				return dataList[i].LikesCount > dataList[j].LikesCount
			})
		}
		//sort.Slice(dataList, func(i, j int) bool {
		//	return dataList[i].LikesCount < dataList[j].LikesCount
		//})

	}
	return dataList, errmsg.SUCCESS
}

func Top() (int, []Topic) {
	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 服务器地址
	})

	// 从 MySQL 中查询数据，并按照字段排序（这里以 ID 字段为例）
	dataList, code := getDataFromMySQL("likes_count", false) // true表示升序，false表示降序
	if code != errmsg.SUCCESS {
		return code, nil
	}
	// 将排序后的数据列表转换成 JSON 格式

	dataJSON, err := json.Marshal(dataList)
	if err != nil {
		return errmsg.ERROR, nil
	}

	// 将 JSON 数据存储到 Redis 的有序集合中，并设置过期时间为五秒
	ctx := context.Background()
	err = client.Set(ctx, "topics", dataJSON, 5*time.Second).Err()
	if err != nil {
		return errmsg.ERROR, nil
	}

	// 从 Redis 中获取排序后的数据并输出
	sortedDataJSON, err := client.Get(ctx, "topics").Result()
	if err != nil {
		return errmsg.ERROR, nil
	}

	var sortedData []Topic
	err = json.Unmarshal([]byte(sortedDataJSON), &sortedData)
	if err != nil {
		return errmsg.ERROR, nil
	}
	// 输出排序后的数据
	return errmsg.SUCCESS, sortedData
}

// GetTopInfo 查询单个话题
func GetTopInfo(id int) (Topic, int) {
	var top Topic
	err = db.Where("id = ?", id).First(&top).Error
	if err != nil {
		return top, errmsg.ERROR_ART_NOT_EXIST
	}
	return top, errmsg.SUCCESS
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
	PointsDe(userid)

	return errmsg.SUCCESS
}

// 点赞
func LikeTopic(topicID, userID uint) int {
	// 检查是否已经点赞过
	var like Like
	_, err := GetTopInfo(int(topicID))
	if err != errmsg.SUCCESS {
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

	return errmsg.SUCCESS
}

// 取消点赞
func Unliketop(topicId, userId uint) int {
	_, code := GetTopInfo(int(topicId))
	if code != errmsg.SUCCESS {
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
	return errmsg.SUCCESS
}

// 添加收藏
func FavoriTopADD(topid, userid uint) int {
	var favo Favorite
	_, err := GetTopInfo(int(topid))
	if err != errmsg.SUCCESS {
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
	return errmsg.SUCCESS
}

// 删除收藏
func DeFavoTop(topid, userid uint) int {
	_, code := GetTopInfo(int(topid))
	if code != errmsg.SUCCESS {
		return errmsg.ERROR_ART_NOT_EXIST
	}
	//删除收藏数
	err = db.Where("topic_id = ? AND user_id = ?", topid, userid).Delete(&Favorite{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
func Favouser(userid uint) (int, []Topic) {
	var favorites []Favorite
	err := db.Where("user_id = ?", userid).Find(&favorites).Error
	if err != nil {
		return errmsg.ERROR, nil
	}
	var tcpids []uint
	for _, favorite := range favorites {
		tcpids = append(tcpids, favorite.TopicID)
	}
	var topics []Topic
	db.Where("id IN (?)", tcpids).Find(&topics)
	return errmsg.SUCCESS, topics
}
