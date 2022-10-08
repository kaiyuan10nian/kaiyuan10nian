### Go实战项目之分页加载

这个功能是前几天就实现了的，忘记提交上去，在开始下一步之前先把这个功能简单介绍一下。

分页加载在日常开发中非常常见，如：

> XXX记录明细
>
> XXX文章列表
>
> XXX好友列表

分页加载好处多多，我们经常会遇到在数据较多的情况下前端一次性加载会比较吃力，处理不好甚至会导致前端的OOM（前几年某知名大型blogAPP就出现过这个BUG）。总之，分页加载肯定是针对这种大量数据报表式的展示最友好的处理方式了。

其实对于后端开发人员来说实现起来也非常简单，核心思想就是：*根据偏移量分多次查询数据即可*。

下面我们分析下对原项目的改动是如何一步一步实现的：

#### 第一步 添加了两个字段分别是page和size

page是页码数，是我们计算偏移量的一个重要元素，之所以设置这个字段是为了方面前端分页使用；

size是每页的数据条目数，也是我们计算偏移量的重要元素，这个字段控制每页返回的数据条目数；

首先在接口中获取这两个字段：

```go
page,_ := strconv.Atoi(ctx.DefaultQuery("page","1"))
size,_ := strconv.Atoi(ctx.DefaultQuery("size","10"))
```

其中设置了默认值，如果前端不给传这两个字段，我们默认为第一页且条目数为10条。

然后先去查一下一共有多少条数据：

```go
	var total int
	c.DB.Model(&model.Article{}).Count(&total)
	if total == 0 {
		response.Fail(ctx,"文章不存在,请创建",nil)
		return
	}
```

查询总条目数的目的就是确定我们库中是否有足够的数据去给前端提供，这里我仅仅判断了0条目状态下的处理，你也可以在这里去添加于页码等不相符的逻辑，比如我们一共只有9条数据，你前端穿了一个page=2并且size=10的值给我，那么可以在这里判断后告诉你：请求数据不足等信息。

#### 第二步 计算偏移量

每次去查数据库的时候，我们先看下前端用户给的是什么值，根据用户需求去查对应的数据：

```go
offset := (page-1)*size
```

比如，用户传入page=1&size=5,那么这里的offset=（1-1）*5=0，偏移量就是0就意味着我们不需要在查数据库是跳过偏移量，直接从第一条开始查够5条后返回给前端即可。

```go
if err := c.DB.Order("id DESC").Offset(offset).Limit(size).Find(&articles).Error;err != nil{
		response.Fail(ctx,"查询失败",nil)
		return
	}
```

上面这段代码实现了偏移量为offset，一共要size条数，并按照id排序的数据给了articles对象。

#### 最后 返回结果给前端

```go 
response.Success(ctx,gin.H{"articles":articles, "page": page,"total": total},"请求成功")
```

直接把数据送回前端，over~

完整代码如下：

```go
//文章列表
func (c ArticleController) List(ctx *gin.Context) {
	page,_ := strconv.Atoi(ctx.DefaultQuery("page","1"))
	size,_ := strconv.Atoi(ctx.DefaultQuery("size","10"))

	var total int
	c.DB.Model(&model.Article{}).Count(&total)
	if total == 0 {
		response.Fail(ctx,"文章不存在,请创建",nil)
		return
	}
	var articles []model.Article
	offset := (page-1)*size
	if err := c.DB.Order("id DESC").Offset(offset).Limit(size).Find(&articles).Error;err != nil{
		response.Fail(ctx,"查询失败",nil)
		return
	}
	//返回结果
	response.Success(ctx,gin.H{"articles":articles, "page": page,"total": total},"请求成功")
	return
}
```

完事！