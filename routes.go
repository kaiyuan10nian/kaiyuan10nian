package main

import (
	"github.com/gin-gonic/gin"
	"kaiyuan10nian/controller"
)

func CollectRoute(r *gin.Engine) *gin.Engine  {
	r.POST("/v1/account/register", controller.Register)

	return r
}