#### 实现图片上传和URL生成

这个功能其实不是很难，对于新手来说比较困难的是对服务器各个路径的掌握。下面先介绍实现图片上传的逻辑，该功能一共分6步来实现：

###### 第一步 获取上传的文件

我们不可能让所有人都可以随意上传东西到我们服务器的，所以在开这个接口的时候肯定是要带token去验证身份的，我们先把接口开出来

```go
r.POST("/v1/upload",common.AuthMiddleware(), controller.Uploads)//上传图片
```

接口开好之后，就是接收上传的文件：

```go
file, err := ctx.FormFile("file")
```

###### 第二步 验证上传文件的合法性

接收到文件之后，我们需要对该文件做一个简单的解析，不是所有的文件都要接收，我们只接收后缀为.jpg .png .jpeg格式的图片：

```go
extName := path.Ext(file.Filename)
allowExtMap := map[string]bool{
   ".jpg":  true,
   ".png":  true,
   ".jpeg": true,
}
if _, ok := allowExtMap[extName]; !ok {
   config.GetLogger().Error(errors.New("文件类型不合法"), "上传错误", false)
   // 返回值
   response.Response(ctx,http.StatusOK,50001,nil,"文件类型不合法")
   return
}
```

如果文件类型不合法则直接拒绝访问。

###### 第三步 创建保存图片的目录

在创建之前，我们首先去初始化一个文件静态地址，Gin自带一个StaticFS函数可以满足我们的需求：

```go
r.StaticFS("/kaiyuan", http.Dir("/opt/server/nginx-1.18/html/kaiyuan"))
```

这个是放在系统初始化的时候去执行的。

回到正题，创建我们保存图片的目录：

```go
currentTime := time.Now().Format("20060102")
// 生成目录文件夹，并错误判断
if err := os.MkdirAll("/opt/server/nginx-1.18/html/kaiyuan/upload/"+currentTime, 0755); err != nil {
   config.GetLogger().Error(err, "上传错误", false)
   // 返回值
   response.Response(ctx,http.StatusOK,50001,nil,"MkdirAll失败")
   return
}
```

为了好管理，我们把文件夹按日期去命名，另外在Linux系统中有一套完善的权限管理机制，为了读取方便我们要给这个目录设置一个0755的权限。另外我们这里用的是MkdirAll函数，它可以级联创建的。

###### 第四部 生成文件名称

Linux系统是不允许同一个目录下有两个同样名称的文件存在的，为了避免这种情况的发生，我们可以把上传来的文件进行一下重命名，按照毫米值去命名将是一个非常不错的方法。

```go
fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)
```

###### 第五步 保存文件

通过以上步骤文件也拿到了，目录也创建好了，命名也处理了，那这一步就是考虑保存了。

```go
saveDir := path.Join("/opt/server/nginx-1.18/html/kaiyuan/upload/"+currentTime, fileUnixName+extName)
err := ctx.SaveUploadedFile(file, saveDir)
if err != nil {
   config.GetLogger().Error(err, "上传错误", false)
   // 返回值
   response.Response(ctx,http.StatusOK,5001,nil,"文件保存失败")
   return
}
```

这一步非常关键，处理不好你就找不到你上传的图片去哪了。

第六步 返回URL

仅仅是上传图片肯定是不够的，因为我们上传就是为了使用，所以一定要把图片的URL给用户返回回去。这里我进行了一下拼接来实现的，如果你有更好的方法请来告诉我哦！

```go
imageurl := strings.Replace(saveDir,"/opt/server/nginx-1.18/html","https://xiaoyin.live",-1)
// 返回值
response.Success(ctx,gin.H{"imageurl":imageurl},"上传成功")
```

至此图片上传的问题就搞定了。

###### 注

这里有个小知识点，如果你域名配置的路径和你图片目录不在一块，那么你就只能通过ip去访问你的图片，比如：132.13.344.12:8080/kaiyuan/upload/20220915/1663222506540589000.jpeg,而如果你图片的目录就在域名配置的路径中，那你就可以直接通过域名去访问了，比如：https://xiaoyin.live/kaiyuan/upload/20220915/1663226910504159272.jpg



ok just this ...  下面把该文件全部内容贴一下，更加详细可前往开源项目github去查阅，或者进群讨论。

```go
package controller

import (
   "errors"
   "github.com/gin-gonic/gin"
   "kaiyuan10nian/config"
   "kaiyuan10nian/response"
   "net/http"
   "os"
   "path"
   "strconv"
   "strings"
   "time"
)

// 上传图片接口
func Uploads(ctx *gin.Context) {
   //1、获取上传的文件
   file, err := ctx.FormFile("file")
   if err == nil {
      //2、获取后缀名 判断类型是否正确 .jpg .png .jpeg
      extName := path.Ext(file.Filename)
      allowExtMap := map[string]bool{
         ".jpg":  true,
         ".png":  true,
         ".jpeg": true,
      }
      if _, ok := allowExtMap[extName]; !ok {
         config.GetLogger().Error(errors.New("文件类型不合法"), "上传错误", false)
         // 返回值
         response.Response(ctx,http.StatusOK,50001,nil,"文件类型不合法")
         return
      }
      //3、创建图片保存目录,linux下需要设置权限（0755可读可写） kaiyuan/upload/image20220915
      currentTime := time.Now().Format("20060102")
      // 生成目录文件夹，并错误判断
      if err := os.MkdirAll("/opt/server/nginx-1.18/html/kaiyuan/upload/"+currentTime, 0755); err != nil {
         config.GetLogger().Error(err, "上传错误", false)
         // 返回值
         response.Response(ctx,http.StatusOK,50001,nil,"MkdirAll失败")
         return
      }
      //4、生成文件名称 1663213319130065587.png
      fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)
      //5、上传文件 kaiyuan/upload/20220915/144325235235.png
      saveDir := path.Join("/opt/server/nginx-1.18/html/kaiyuan/upload/"+currentTime, fileUnixName+extName)
      err := ctx.SaveUploadedFile(file, saveDir)
      if err != nil {
         config.GetLogger().Error(err, "上传错误", false)
         // 返回值
         response.Response(ctx,http.StatusOK,5001,nil,"文件保存失败")
         return
      }
      imageurl := strings.Replace(saveDir,"/opt/server/nginx-1.18/html","https://xiaoyin.live",-1)
      // 返回值
      response.Success(ctx,gin.H{"imageurl":imageurl},"上传成功")
      return
   }
}
```