package user

import (
	"bbs-go/model"
	"bbs-go/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpLoad(c *gin.Context) {
	id := c.GetUint("userid")
	var profile model.Profile
	file, flileHeader, _ := c.Request.FormFile("file")
	fileSize := flileHeader.Size
	url, code := model.UpLoadFile(file, fileSize)
	profile.Avatar = url
	if code == errmsg.SUCCSE {
		code = model.UpdateProfile(int(id), &profile)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"url":     url,
	})
}
