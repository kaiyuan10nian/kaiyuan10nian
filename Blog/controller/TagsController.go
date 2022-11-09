package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"kaiyuan10nian/common"
	"kaiyuan10nian/model"
	"kaiyuan10nian/response"
	"kaiyuan10nian/vo"
	"net/http"
	"strconv"
)

type ITagsController interface {
	RestController
}

type TagsController struct {
	DB *gorm.DB
}



func NewTagController() ITagsController {
	db := common.GetDB()

	return TagsController{DB:db}
}
//创建标签
func (c TagsController) Create(ctx *gin.Context) {
	var requestTag vo.CreateTagRequest
	//先去取穿上来的数据并进行验证
	if err := ctx.ShouldBind(&requestTag); err != nil {
		response.Fail(ctx, "数据验证错误，标签名称必填", nil)
		return
	}
	//通过验证后再去数据库表中创建对应记录
	tag := model.Tags{TagName:requestTag.TagName}
	if err := c.DB.Create(&tag).Error;err != nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,err.Error())
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"tag": tag}, "")
	return
}
//更新标签
func (c TagsController) Update(ctx *gin.Context) {
	var requestTag vo.CreateTagRequest
	//先去取穿上来的数据并进行验证
	if err := ctx.ShouldBind(&requestTag); err != nil {
		response.Fail(ctx, "数据验证错误，标签名称必填", nil)
		return
	}

	tagId,_ := strconv.Atoi(ctx.Params.ByName("id"))
	//然后查一下要修改的这个标签记录存在不存在
	var updateTag model.Tags
	if c.DB.First(&updateTag,tagId).RecordNotFound(){
		response.Fail(ctx,"标签不存在",nil)
		return
	}
	//存在的话就去修改
	if err := c.DB.Model(&updateTag).Update("tag_name",requestTag.TagName).Error;err != nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,err.Error())
		return
	}
	//返回修改结果
	response.Success(ctx,gin.H{"tag":updateTag},"修改成功")
	return
}
//标签详情
func (c TagsController) Show(ctx *gin.Context) {
	tagId,_ := strconv.Atoi(ctx.Params.ByName("id"))
	//根据标签ID直接去表中查对应标签
	var tag model.Tags
	if c.DB.First(&tag,tagId).RecordNotFound(){
		response.Fail(ctx,"标签不存在",nil)
		return
	}
	//返回结果
	response.Success(ctx,gin.H{"tag":tag},"")
	return
}
//标签列表
func (c TagsController) List(ctx *gin.Context) {
	var tags []model.Tags
	//直接去查询所有标签
	c.DB.Find(&tags)
	var total int
	c.DB.Model(&model.Tags{}).Count(&total)
	if total == 0 {
		response.Fail(ctx,"标签不存在,请创建",nil)
		return
	}
	//返回结果
	response.Success(ctx,gin.H{"tags":tags, "total": total},"")
	return
}
//删除标签
func (c TagsController) Delete(ctx *gin.Context) {
	tagId,_ := strconv.Atoi(ctx.Params.ByName("id"))

	if err := c.DB.Delete(model.Tags{},tagId).Error;err != nil{
		response.Fail(ctx,"删除失败请重试",nil)
		return
	}

	response.Success(ctx,nil,"删除成功")
	return
}
