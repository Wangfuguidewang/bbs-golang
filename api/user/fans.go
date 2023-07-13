package user

import (
	"bbs-go/model"
	"bbs-go/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FansAdd(c *gin.Context) {
	follid, _ := strconv.Atoi(c.Param("follid"))
	userid := c.GetUint("userid")
	code := model.FollwAdd(uint(follid), userid)
	c.JSON(http.StatusOK, gin.H{
		"staus":   code,
		"message": errmsg.GetErrMsg(code),
	})
}
func Unfans(c *gin.Context) {
	follid, _ := strconv.Atoi(c.Param("follid"))
	userid := c.GetUint("userid")
	code := model.UnFollw(uint(follid), userid)
	c.JSON(http.StatusOK, gin.H{
		"staus":   code,
		"message": errmsg.GetErrMsg(code),
	})
}
