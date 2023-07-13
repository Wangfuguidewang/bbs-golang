package user

import (
	"bbs-go/model"
	"bbs-go/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Createtopic(c *gin.Context) {
	//tokenString := r.Header.Get("Authorization")
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	return []byte("you-secret-key"), nil
	//})
	//if err != nil || !token.Valid {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}
	//
	//// 从 token 中提取用户信息，例如用户 ID
	//claims := token.Claims.(jwt.MapClaims)
	//userID := uint(claims["userID"].(float64))
	id := c.GetUint("userid")
	var data model.Topic
	//data.UserID = userID
	data.UserID = id
	_ = c.ShouldBindJSON(&data)
	code = model.CreateTopic(&data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetCateTop 查询userid下的所有话题
func GetCateTop(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	id := c.GetUint("userid")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	data, code, total := model.GetCaTopic(int(id), pageSize, pageNum)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// GetTopInfo 查询单个话题信息
func GetTopInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetTopInfo(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
} // EditTop 编辑话题
func EditTop(c *gin.Context) {
	var data model.Topic
	userid := c.GetUint("userid")
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := model.EditTopic(id, userid, &data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

// DeleteTop 删除话题
func DeleteTop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userid := c.GetUint("userid")
	code = model.DeleteTop(id, userid)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}
