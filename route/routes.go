package route

import (
	"github.com/gin-gonic/gin"
	"kaiyuan10nian/common"
	"kaiyuan10nian/config"
	"kaiyuan10nian/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine  {
	r.Use(config.LoggerToFile())//添加日志记录
	r.POST("/v1/account/register", controller.Register)//注册
	r.POST("/v1/account/login", controller.Login)//登录
	r.GET("/v1/account/info", common.AuthMiddleware(),controller.Info)//用户信息
	r.GET("/v1/account/invite", common.AuthMiddleware(),controller.Invite)//邀请码生成
	r.POST("/v1/upload",common.AuthMiddleware(), controller.Uploads)//上传图片
	tagsRoutes := r.Group("/v1/tags")//标签
	tagsRoutes.Use(common.AuthMiddleware())
	tagController := controller.NewTagController()
	tagsRoutes.POST("", tagController.Create)
	tagsRoutes.PUT("/:id", tagController.Update) //替换
	tagsRoutes.GET("/:id", tagController.Show)
	tagsRoutes.DELETE("/:id", tagController.Delete)
	tagsRoutes.GET("/list", tagController.List)
	articleRoutes := r.Group("/v1/article")//文章
	articleRoutes.Use(common.AuthMiddleware())
	articleController := controller.NewArticleController()
	articleRoutes.POST("", articleController.Create)
	articleRoutes.PUT("/:id", articleController.Update) //替换
	articleRoutes.GET("/:id", articleController.Show)
	articleRoutes.DELETE("/:id", articleController.Delete)
	articleRoutes.GET("/list", articleController.List)
	test := r.Group("/test")
	{
		test.GET("/hello", func(context *gin.Context) {
			msg := "ok"
			context.JSON(200, msg)
		})
	}

	return r
}