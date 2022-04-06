## <a id="8">JWT-Go 实现登录Token发放和验证</a>

第一阶段的知识点啃的差不多了，这周继续往下进行。因为我们目前要实现的是一个BLOG系统，那么我们肯定要考虑安全问题。不能说任何一个人上来就能随便发布博文，那我们这个blog系统最后就广告、垃圾横行了。所以我们只允许注册了的用户使用我们的系统，没有经过注册的坚决拒绝其使用。

在系统中怎么判断一个人是否是我们的用户应该怎么做呢？那就是用户拿自己的账号和密码给系统进行验证，系统看看他是不是已经注册的用户，如果账号和密码都能匹配上就允许访问，如果匹配不上就拒绝其访问。

但是在实际项目中，我们不可能在用户每次访问的时候都跟用户去要账号和密码，就算你这么去要也不安全呀。所以这个时候我们就需要采取一种加密手段。通过这种加密手段去处理用户的账号和密码，然后在用户每次访问我们的时候带上这一穿加密字符串就行了。

这样既解决了认证的问题还解决了安全的问题，这就是Token，在Go语言中有一个库已经解决了我们这个问题他就是：jwt-go.

安装go-jwt：go get -u github.com/dgrijalva/jwt-go

在开始写代码之前我们先回顾一下流程：

1）用户注册账号设置密码，成功后我们把其信息存入我们的库中。

2）用户登录，服务端先校验用户传入的账号和密码，准确无误后生成一个token给用户。

3）用户访问需要权限的接口时携带这个token即可。

下面开始编程，一探jwt-go的用法。

#### 1、指定加密秘钥

```go
var jwtKey = []byte("kai_yuan_shi_nian")
```

这就相当于是一把钥匙，自己保存好，造锁和开锁都是依托这把钥匙进行的。

#### 2、创建Claims结构体

```go
type Claims struct {
	UserId uint
	jwt.StandardClaims
}
```

这个结构体就是用来保存信息的，需要内嵌jwt.StandardClaims，这些信息会被保存在我们生成好的token当中。

#### 3、生成token

```go
func ReleaseToken(user model.User) (string,error){
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId : user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),//设置这个token的有效期
			IssuedAt: time.Now().Unix(),//发放时间
			Issuer: "kaiyuanshinian.tech",//发行方
			Subject: "user token",//主题
		},
	}
	//使用指定的签名方式创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
  //使用上面指定的钥匙(secret)签名并获取完整的签名后的字符串
	tokenString,err := token.SignedString(jwtKey)

	if err != nil{
		return "",err
	}

	return tokenString,nil
}
```

在上面这段代码中我们分别设置了有效期和发布方信息。除了上面代码中用到的有效期、签发时间、签发人信息外，还有生效时间(NotBefore)、受众(Audience)、编号(JWTID)等等信息看需求你可以自己添加。

#### 4、解析token

```go
func ParseToken(tokenString string)(*jwt.Token,*Claims,error){
	claims := &Claims{}
	//用于解析鉴权声明
	token,err := jwt.ParseWithClaims(tokenString,claims,func(token *jwt.Token)(i interface{},err error){
		return jwtKey,nil
	})

	return token,claims,err
}
```

解析token就是在用户访问我们的时候，我们系统去解析他所携带的token，去验证它是否是我们正确的用户。我们可以直接根据token获取到它所携带的用户信息（上面的结构体）

#### 5、编写路由

上面完成后，我们开始写我们的业务功能，前面我们已经实现注册功能了，这里就不多说了。首先我们去实现登录功能，直接上代码分析：

```go
r.POST("/v1/account/login", controller.Login)
```

先写路由，指向Login去接收处理，那么我们再看看系统收到该接口访问时的处理。

其处理顺序是：

1）先获取用户访问接口是携带的参数

2）拿这些参数去校验我们的库是否准确（手机号、密码等信息）

3）准确的话发放token并返回登录成功

```go
//登录
func Login(ctx *gin.Context){
	DB := common.GetDB()
	//获取参数
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	//数据验证
	if len(mobile) != 11 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	DB.Where("mobile = ?",mobile).First(&user)
	if user.ID == 0 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password));err != nil{
		response.Response(ctx,http.StatusUnprocessableEntity,400,nil,"密码错误")
		return
	}
	//发放token
	token,err := common.ReleaseToken(user)
	if err != nil{
		response.Response(ctx,http.StatusUnprocessableEntity,500,nil,"系统异常")
		log.Printf("token generate error : %v",err)
		return
	}
	//返回结果
	fmt.Println(token)
	response.Success(ctx,gin.H{"token":token},"登录成功")
}
```

在上面的代码中我们就用到了前面实现的token发放模块（ReleaseToken）。

到此，登录功能就已经实现了，那么这个token怎么使用呢？我们继续往下去实现一个可以获取用户信息的接口。

#### 6、token验证及使用

直接先上路由

```go
r.GET("/v1/account/info", common.AuthMiddleware(),controller.Info)
```

各位发现没有，这次注册路由的时候明显与前面两次不同了，多了一个common.AuthMiddleware()。它就是验证token必不可少的一环。

其AuthMiddleware()函数的主要内容如下：

```go 
func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		//先从header中获取token
		tokenString := ctx.GetHeader("Authorization")

		//然后再去验证token不为空和它的类型
		if tokenString == "" || !strings.HasPrefix(tokenString,"Bearer"){
			ctx.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}
		//抛去前面的7个字节不要（其主要内容就是一个bearer类型声明）
    //token示例： ~~Bearer ~~ eyJhbGciOiJIUzI1NiIsInR...Ndafg，前面7位正好是：Bearer <-这里还有一个空格哦

		tokenString = tokenString[7:]

		token,claims,err := ParseToken(tokenString)

		if err != nil || !token.Valid{
			ctx.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"权限不足"})
			ctx.Abort()
			return
		}
		//通过验证后获取claims中的userID
		userId := claims.UserId
		DB := GetDB()
		var user model.User
		DB.First(&user,userId)

		//检查用户是否存在
		if user.ID == 0{
			ctx.JSON(http.StatusUnauthorized,gin.H{"code":401,"msg":"用户不存在"})
			ctx.Abort()
			return
		}

		//如果用户存在 将user信息存入上下文
		ctx.Set("user",user)
		ctx.Next()
	}
}

```

这里主要做了几件事情：

1）校验token的有效性，有效go on ,无效 stop it.

2）校验完毕后取出token中的claims进行解析，根据上面我们的结构体可知，我们可以得到用户的userID。

3）拿该userID去我们的库中查询是否存在，存在go on ,不存在 stop it.

完事。

后面的就不用多说了，看下我们调用v1/account/info接口给我们返回了什么数据吧。

```json
{
    "code": 200,
    "data": {
        "user": {
            "name": "张三",
            "telephone": "13523422342"
        }
    },
    "msg": "请求成功"
}
```

目前我们仅仅用到了jwt-go的核心部分，其还有很多可扩展功能有待我们去开发学习，本节就先到这里吧。

ok just it ,see you next...

