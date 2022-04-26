package route

import (
	"github.com/gin-gonic/gin"
	"kaiyuan10nian/common"
	"kaiyuan10nian/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine  {
	r.POST("/v1/account/register", controller.Register)
	r.POST("/v1/account/login", controller.Login)
	r.GET("/v1/account/info", common.AuthMiddleware(),controller.Info)
	r.GET("/v1/account/invite", common.AuthMiddleware(),controller.Invite)
	tagsRoutes := r.Group("/v1/tags")
	tagsRoutes.Use(common.AuthMiddleware())
	tagController := controller.NewTagController()
	tagsRoutes.POST("", tagController.Create)
	tagsRoutes.PUT("/:id", tagController.Update) //替换
	tagsRoutes.GET("/:id", tagController.Show)
	tagsRoutes.DELETE("/:id", tagController.Delete)
	tagsRoutes.GET("/list", tagController.List)
	return r
}