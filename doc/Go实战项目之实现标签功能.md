今天在考虑怎么实现文章分组管理的逻辑，想实现类似于下图这样的一个功能：

![WechatIMG290.png](https://s2.loli.net/2022/04/12/UA7Jm6g1otTDPKG.png)

左边是文章类型分组，右边是该分组下面的文章。本来想的是创建一个分组然后每个文章都归属到一个组别去，类似一个多对一的关系。但是后面在做的过程中发现这个思路想的太简单了，而且对后期的扩展特别不友好。

比如，我如果后面想扩展一个标签或者关键字的功能的话，这个分组就非常鸡肋了。所以直接放弃了之前的想法，仔细想想的话，分组无非就是按照文章的不同属性去给划分了类别，那么这个不同的属性其实就是标签了，一篇文章可以有很多标签，那么在分组的过程中我们就可以按照标签去实现。

那么思路明确了，下面就开始我们标签功能的开发：

1）创建数据库  两个表 一个是标签表，一个是标签和文章的关系表

2）创建路由，针对标签的增删查改

这一次的功能非常简单就这么多，但是在我实现的过程中有几个知识点我认为应该分享一下：

1）gin的binding and shouldbind
2）gorm的软删除
3）go 中 interface的使用
4）gin route 中group分组使用

我们直接开始去实现功能，到涉及知识点的地方我重点标记一下。

第一步，创建两个model，分别是标签和标签与文章关系

首先是标签的model，其实就是一个标签ID和标签名称就够了

```go
type Tags struct {
	gorm.Model
	TagName  string      `json:"tagname" gorm:"not null"`
}
```

标签和文章关系表中也不复杂，一个关系ID，一个标签ID，一个文章ID足以

```go
type TagMapArticle struct {
	gorm.Model
	ArticleID         uuid.UUID `json:"article_id" gorm:"not null"`
	TagID  uint      `json:"tag_id" gorm:"not null"`
}
```

好的，完成以上部分就ok了，下面就是启用这两个model，直接在数据库初始化的时候去创建这两张表

```go
	db.AutoMigrate(&model.Tags{})
	db.AutoMigrate(&model.TagMapArticle{})
```

第一步就算完成了

第二步，开功能接口

标签的管理肯定离不开针对标签的增删查改四个维度，除此之外在需求习惯中我们一般会再额外实现一个展示标签列表的功能

按照之前的逻辑我们是不是应该如下面代码所示去实现呀？

```go
r.POST("/v1/tags/create", controller.Create)
r.POST("/v1/tags/update", controller.Update)
r.GET("/v1/tags/show", controller.Show)
r.DELETE("/v1/tags/delete", controller.Delete)
r.GET("/v1/tags/list", controller.List)
```

这么写可不可以，当然可以了。但是，对于我们后期维护以及扩展也是非常不友好的，因为在实际项目中我们都是分模块的，比如说安全模块、账户管理模块、文章管理模块、标签管理模块等等，另一方面既然是分模块开发，那么我们就不能去把所有处理都去写进同一个controller中。

在大型项目中我们一般都会采取一些架构设计，比如什么MVC、MVVC、MVP等等，如果你上面那种写法很难去实现了。

所以这里GIN给我们提供了一个非常好用的功能：**路由分组**

路由分组就是把同一个模块的或者同一个版本的去放到一个组别中去，然后统一对组内的路由进行管理，如果我们要用路由分组去实现上面的几个路由应该怎么写呢？

```go
	tagsRoutes := r.Group("/v1/tags")
	tagController := controller.NewTagController()
	tagsRoutes.POST("", tagController.Create)
	tagsRoutes.PUT("/:id", tagController.Update) 
	tagsRoutes.GET("/:id", tagController.Show)
	tagsRoutes.DELETE("/:id", tagController.Delete)
	tagsRoutes.GET("/list", tagController.List)
```

针对Tag的管理路由全部放进了同一个组内，统一设置了前缀“/v1/tags”，后面路由可以再根据自己的功能去做具体区分。

好了，接口功能实现了，但是我们在项目中这个标签的管理肯定不能公开的去放出去，所以借鉴前面讲过的给这几个接口加上权限：

```go
tagsRoutes.Use(common.AuthMiddleware())
```

这样，路由就完成了。

路由完成并不是就完事了，我们要针对不同的路由去实现对应的逻辑，由于我们实现了分组路由管理，那么针对标签管理模块我独立创建了一个controller去专门管理TagsController.go，具体内容我们在代码中标记的非常清楚，下面就是功能逻辑的实现

```go

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
```

这几行代码中的增删查改逻辑并不复杂，所以我觉得基本大家都看的懂，这里涉及的知识点就是：**interface多态**

很多新手不喜欢用interface，总感觉我实例化一下直接就去调相关函数了整这些花里胡哨的干什么？因为我之前是JAVA出身在最开始就是这种感觉，但是随着我项目经验的日积月累就再次验证了真香定律。

给大家举个例子：

你走在大街上看到一家店的招牌上写着“麦当劳”，那你不用进店就知道它里面卖的是汉堡、薯条、可乐、鸡翅等。而如果你看到一家店的招牌上写着“李二厨”，请问你知道他店里是都有什么菜吗？

”麦当劳“这三个字就类似于代码中的接口，**不管在哪个城市哪条街哪家店只要挂了**这个招牌就相当于是实现了”汉堡、薯条、可乐、鸡翅“这几个函数。

这其实就是**多态**，更多概念这里就不去细讲了，有兴趣可以去自己查查。

它的优点你要记住，后面自己去验证一下，每一个优点都能举100个例子证明，**这是开发人员最重要的知识点之一**：

1）可替换性

2）可扩充性

3）接口性

4）灵活性

5）简化性

上面的多态是一个知识点，另外一个你在代码中应该发现了：

```go 
//先去取穿上来的数据并进行验证
	if err := ctx.ShouldBind(&requestTag); err != nil {
		response.Fail(ctx, "数据验证错误，标签名称必填", nil)
		return
	}
```

怎么直接通过ctx.ShouldBind去获取参数了呢？在前面的项目开发中我们是怎么获取接口传上来的参数的？

```go
	name := ctx.PostForm("name")
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	recommender := ctx.PostForm("recommender")
```

是不是不一样？这就是**Gin框架的参数绑定**。

在实战项目中我们避免不了要写很多接口，肯定会设计参数的传递，无论是path/query/string还是body都是避免不了的事情，那么你是否遇到过如下问题：

1）我写的两个功能接口，第二个接口仅仅比第一个接口多了一个ID的参数，那我能不能复用代码？（不考虑复用的不是一个合格的程序员）

2）model中的Struct能不能跟我的参数绑定？

不用说，你肯定遇到过，没遇到只能说你项目做的还不够多。那么怎么解决？

**Gin参数绑定**来帮你解决！！！

在Gin中，为我们提供了一些列的binding函数让我们可以把这些参数绑定到一个对象中，还可以通过struct tag 来对参数进行校验。

具体知识点内容请自行查询，相信通过自己努力获得的知识点会记得更深，这里我只是告诉你什么东西可以解决什么问题。

接下来就是设计到的第四个知识点

![1650697602351.jpg](https://s2.loli.net/2022/04/23/4Jdu9Yy5g8hQIzc.jpg)

看上面数据库中表记录能看出什么不？

是的，第一条有一个删除时间，其他是没有的。这就是**GORM中的软删除**。

如果你的model中包含了gorm.DeletedAt字段（包含在gorm.Model）,就将自动获得软删除能力！

当你在调用Delete时，记录不会从数据库中删除，但GORM会将DeletedAt字段的值设置为当前时间，并且在你使用正常的Query方法查找数据时将不会返回该记录。

ok just this see you next...