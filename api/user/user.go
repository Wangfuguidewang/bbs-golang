package user

import (
	"bbs-go/middleware"
	"bbs-go/model"
	"bbs-go/utils/errmsg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var code int

// 添加用户
func AddUser(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println(data)
		return
	}
	code = model.CheckUser(data.Username, data.Email)
	if code == errmsg.SUCCSE {
		//返回值正确 表示没有用户 则可添加
		data.LastLoginTime = time.Now()
		model.CreateUser(&data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		code = errmsg.ERROR_USERNAME_USED

	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})

}

// 登录
func Login(c *gin.Context) {
	var data model.User
	var token string
	var code int
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("err ==", err)
		return
	}
	code = model.CheckLOgin(data.Username, data.Password)

	if code == errmsg.SUCCSE {
		code, userid := model.CheckUserid(data.Username)
		if code == 1001 {
			fmt.Println("code ==", code)
			return
		}
		data.LastLoginTime = time.Now()
		token, code = middleware.SetToken(userid, data.Password)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"token":   token,
	})
}
