package user

import (
	"bbs-go/model"
	"bbs-go/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProfile(c *gin.Context) {
	id := c.GetUint("userid")
	data, code := model.Getprofile(int(id))
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": errmsg.GetErrMsg(code),
	})
}
func UpdataProfile(c *gin.Context) {
	var data model.Profile
	data.ID = c.GetUint("userid")
	_ = c.ShouldBindJSON(&data)
	code := model.CheckUser(data.Username, data.Email)
	if code == errmsg.ERROR_EMAIL_USED {
		code, email := model.Checkenail(data.Email)
		if email == data.Email {
			code = code
		}
	}
	code = model.UpdateProfile(int(data.ID), &data)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})

}
