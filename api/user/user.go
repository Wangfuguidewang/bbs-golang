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
var verificationCodes string

// 添加用户
func AddUser(c *gin.Context) {
	var data model.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println(data)
		return
	}
	code = model.CheckUser(data.UserName, data.Email)

	if code == errmsg.SUCCESS {
		//返回值正确 表示没有用户 则可添加
		data.LastLoginTime = time.Now()
		code = model.CreateUser(&data)

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
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println("err ==", err)
		return
	}
	code, userInfo := model.CheckLogin(data.UserName, data.Password)

	if code == errmsg.SUCCESS {
		data.LastLoginTime = time.Now()
		token, code = middleware.SetToken(userInfo.ID, data.Password)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"token":   token,
	})
}

// 邮箱验证接口
func verifyEmailHandler(c *gin.Context) {
	// 获取用户提交的验证码
	code := c.PostForm("code")

	fmt.Println(code)
}
