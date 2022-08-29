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
	tagsRoutes := r.Group("/v1/tags")//标签
	tagsRoutes.Use(common.AuthMiddleware())
	tagController := controller.NewTagController()
	tagsRoutes.POST("", tagController.Create)
	tagsRoutes.PUT("/:id", tagController.Update) //替换
	tagsRoutes.GET("/:id", tagController.Show)
	tagsRoutes.DELETE("/:id", tagController.Delete)
	tagsRoutes.GET("/list", tagController.List)
	test := r.Group("/test")
	{
		test.GET("/hello", func(context *gin.Context) {
			msg := "ok"
			context.JSON(200, msg)
		})
	}

	return r
}