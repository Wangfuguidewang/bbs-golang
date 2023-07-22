package routes

import (
	"bbs-go/api/user"
	"bbs-go/middleware"
	"bbs-go/utils"
	"github.com/gin-gonic/gin"
)

// 接口
func InitRouter() {
	go func() {
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
			auth.GET("admin/getcate", user.GetCateTop)      //查询用户所有话题
			auth.GET("admin/getinfo/:id", user.GetTopInfo)  //查询单个话题
			auth.PUT("admin/edit/:id", user.EditTop)        //编辑话题
			auth.DELETE("admin/delect/:id", user.DeleteTop) //删除话题
			//评论模块
			auth.POST("addcomment/:topicid", user.AddComment) //添加评论
			//点赞
			auth.POST("likeadd/:topicid", user.Liketopic)  //添加点赞
			auth.DELETE("unlike/:topicid", user.Unliketop) //删除点赞
			//收藏
			auth.POST("favAdd/:topicid", user.FavotopAdd) //添加收藏
			auth.DELETE("favDe/:topicid", user.DeFavoTop) //删除收藏
			auth.GET("favouser", user.FavoUser)           //查看用户下的所有搜藏

			//关注粉丝
			auth.POST("folladd/:follid", user.FansAdd) //添加关注 粉丝+1
			auth.DELETE("Defool/:follid", user.Unfans) //取消关注 粉丝-1
		}
		router := r.Group("/api")
		{
			//用户信息模块
			router.POST("user/add", user.AddUser)
			router.POST("user/email", user.Emailadd)       //发送验证码
			router.POST("user/email/code", user.Emailcode) //验证验证码
			//登陆控制模块
			router.POST("user/login", user.Login)
			//查询话题列表
			router.GET("topic", user.GetTop)
		}
		r.Run(utils.HttpPort)
	}()

}
