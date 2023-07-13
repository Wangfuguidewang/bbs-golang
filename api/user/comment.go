package user

import (
	"bbs-go/model"
	"bbs-go/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddComment 新增评论
func AddComment(c *gin.Context) {
	topicid, _ := strconv.Atoi(c.Param("topicid"))
	userid := c.GetUint("userid")

	var data model.Comment
	_ = c.ShouldBindJSON(&data)
	data.UserID = userid
	_, code := model.GetTopInfo(topicid)
	if code == 200 {
		data.TopicID = uint(topicid)
		code = model.AddComment(&data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}
