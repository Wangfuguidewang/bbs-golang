package user

import (
	"bbs-go/model"
	"bbs-go/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 发送验证码
func Emailadd(c *gin.Context) {
	var data model.User
	c.ShouldBindJSON(&data)
	code := model.GenerateVerificationCode(data.Email)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": errmsg.GetErrMsg(code),
	})
}

// 验证验证码
func Emailcode(c *gin.Context) {
	var data model.User
	c.ShouldBindJSON(&data)
	code := model.VerifyCode(data.Email, data.Code)
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": errmsg.GetErrMsg(code),
	})
}
