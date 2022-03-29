package main

import (
	"github.com/gin-gonic/gin"
	"kaiyuan10nian/common"
	"kaiyuan10nian/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine  {
	r.POST("/v1/account/register", controller.Register)
	r.POST("/v1/account/login", controller.Login)
	r.GET("/v1/account/info", common.AuthMiddleware(),controller.Info)
	return r
}