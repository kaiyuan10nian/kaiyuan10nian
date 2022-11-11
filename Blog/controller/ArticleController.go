package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"kaiyuan10nian/Blog/common"
	"kaiyuan10nian/Blog/model"
	"kaiyuan10nian/Blog/response"
	"kaiyuan10nian/Blog/vo"
	"net/http"
	"strconv"
)

type IArticleController interface {
	RestController
}
type ArticleController struct {
	DB *gorm.DB
}

func NewArticleController() IArticleController {
	db := common.GetDB()

	return ArticleController{DB: db}
}

//创建文章
func (c ArticleController) Create(ctx *gin.Context) {
	var requestArticle vo.CreateArticleRequest
	//先去取穿上来的数据并进行验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		response.Fail(ctx, "数据验证错误，文章内容必填", nil)
		return
	}
	// 获取登录用户
	user, _ := ctx.Get("user")
	// 创建post
	articleContent := model.Article{
		UserId:  user.(model.User).ID,
		Title:   requestArticle.Title,
		Content: requestArticle.Content,
	}
	if err := c.DB.Create(&articleContent).Error; err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"articleContent": articleContent}, "文章发表成功")
	return
}

//更新文章
func (c ArticleController) Update(ctx *gin.Context) {
	var requestArticle vo.CreateArticleRequest
	//先去取穿上来的数据并进行验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		response.Fail(ctx, "数据验证错误，文章内容必填", nil)
		return
	}

	articleId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	//然后查一下要修改的这个标签记录存在不存在
	var updateArticle model.Article
	if c.DB.First(&updateArticle, articleId).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}
	//存在的话就去修改
	if err := c.DB.Model(&updateArticle).Update("title", requestArticle.Title).Error; err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	if err := c.DB.Model(&updateArticle).Update("content", requestArticle.Content).Error; err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	//返回修改结果
	response.Success(ctx, gin.H{"title": updateArticle}, "修改成功")
	return
}

//文章详情
func (c ArticleController) Show(ctx *gin.Context) {
	articleId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	//根据文章ID直接去表中查对应文章
	var article model.Article
	if c.DB.First(&article, articleId).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"article": article}, "")
	return
}

//文章列表
func (c ArticleController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	var total int
	c.DB.Model(&model.Article{}).Count(&total)
	if total == 0 {
		response.Fail(ctx, "文章不存在,请创建", nil)
		return
	}
	var articles []model.Article
	offset := (page - 1) * size
	if err := c.DB.Order("id DESC").Offset(offset).Limit(size).Find(&articles).Error; err != nil {
		response.Fail(ctx, "查询失败", nil)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"articles": articles, "page": page, "total": total}, "请求成功")
	return
}

//删除文章
func (c ArticleController) Delete(ctx *gin.Context) {
	articleId, _ := strconv.Atoi(ctx.Params.ByName("id"))

	if err := c.DB.Delete(model.Article{}, articleId).Error; err != nil {
		response.Fail(ctx, "删除失败请重试", nil)
		return
	}

	response.Success(ctx, nil, "删除成功")
	return
}
