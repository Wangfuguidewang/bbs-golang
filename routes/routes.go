package routes

import (
	"bbs-go/api/user"
	"bbs-go/middleware"
	"bbs-go/utils"
	"github.com/gin-gonic/gin"
)

// 接口
func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(middleware.Loggoer())
	r.Use(gin.Recovery())
	//r := gin.Default()
	auth := r.Group("/api")
	auth.Use(middleware.JwtToken())
	{
		//上传文件接口
		auth.POST("upload", user.UpLoad)
		//更新个人设置
		auth.GET("admin/profile", user.GetProfile) //查询
		auth.PUT("profile", user.UpdataProfile)    //展示
		//话题模块
		auth.POST("admin/topic", user.Createtopic)      //新增
		auth.GET("admin/getcate", user.GetCateTop)      //查询用户所有文章
		auth.GET("admin/getinfo/:id", user.GetTopInfo)  //查询单个话题
		auth.PUT("admin/edit/:id", user.EditTop)        //编辑话题
		auth.DELETE("admin/delect/:id", user.DeleteTop) //删除话题
		//评论模块
		auth.POST("addcomment/:topicid", user.AddComment) //添加评论
	}
	router := r.Group("/api")
	{
		//用户信息模块
		router.POST("user/add", user.AddUser)
		//登陆控制模块
		router.POST("user/login", user.Login)
	}
	r.Run(utils.HttpPort)
}
