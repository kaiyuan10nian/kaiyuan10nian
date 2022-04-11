## Go实战项目-给项目添加一个邀请码的功能

前面总结里面说过，我要给我们的blog项目添加一个邀请码的功能，因为这个项目做完之后我要部署出去的，虽说不靠这个东西去盈利，但是实战要有实战的样子至少把一个项目从0做到1的过程还是要有的。这里就涉及到一个很严重的问题：怎么控制内容？（在国内不可能不控制内容的）

- 人工去审核（我没那么多时间）
- 接入第三方内容审核（要花钱的）
- 智能审核（没这个技术~后面也许可以考虑在开源十年中去实现它）

所以没有办法了，既然我控制不了内容，那我就从源头控制：直接控制输出内容的人。

所以我想到的办法就是在注册的时候做一个限制，不能让每个人都可以注册成功，只能让那些得到我们这个圈子成员信任的人才可以注册体验。那么实现的逻辑要提前去思考了：

- 邀请码怎么来的？

​	邀请码是系统生成的，由已注册会员手动生成，没给会员允许生成5个邀请码，每个邀请码只能使用1次。

- 邀请码在哪使用？

​	邀请码在新会员注册的时候使用（为必填项）

- 邀请码怎么追溯？

​	因为邀请码的生成是已注册会员手动生成的，所以我会在库里面去拉一张表单独去记录，如果有新会员使用了邀请码那么我也会记录下来，后面若想去追溯就很简单了

根据上面的思路，我们一步一步来实现其功能：

1、先写接口路由：

```go
r.GET("/v1/account/invite", common.AuthMiddleware(),controller.Invite)
```

因为邀请码只能是已注册会员去生成，所以必须加上common.AuthMiddleware()权限限制。

2、实现邀请码的生成：

```go
func RandomString(n int) string {
	var letters = []byte("23456789qwertyupkjhgfdsazxcvbnmMNBVCXZASDFGHJKPOUYTREWQ")
	result := make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
```

这个其实就是一个生成随机值的函数，其中我把'0'、'1'、'i'、'l'、'I'、'L'等不容易辨认的字符都去掉了。

3、实现邀请码接口功能：

```go
//创建用户邀请码
func Invite(ctx *gin.Context){
	udata ,_:= ctx.Get("user")//首先通过token去获取用户信息
	dto := dto.ToUserDto(udata.(model.User))
	DB := common.GetDB()
	var inviteCodes []model.InviteCode
	DB.Where("userid = ?",dto.ID).Find(&inviteCodes)//根据用户信息去查询这个用户有多少个邀请码
	fmt.Println(len(inviteCodes))
	if len(inviteCodes) < 5{//如果不够5个则允许他继续生成邀请码
		code := RandomString(5)
		newCode := model.InviteCode{
			Userid: dto.ID,
			Code: code,
			Status: 0,
		}
    //生成成功后把它记录进库表里面
		if err := DB.Create(&newCode).Error;err != nil{
			response.Response(ctx,http.StatusInternalServerError,500,nil,err.Error())
			return
		}
		response.Response(ctx,http.StatusOK,6000,gin.H{"inviteCode":code},"邀请码生成成功")
	}else{
		response.Response(ctx,http.StatusOK,60001,nil,"每个人只可以拥有5个邀请码")
		return
	}

}
```

上述提到了邀请码的创建，我们创建了一个表invitecode专门用来存储（下文有model代码）

4、去修改注册功能

原本我们的注册只需要手机号并设置密码即可，但是我们要控制注册了，那么要加几个验证条件：

1）推荐人是否存在？

2）邀请码是否存在？

3）邀请码是否和推荐人匹配？

4）邀请码是否被人使用过了？

这几个条件判断之后，如果都通过那么给予放行，放行后还要做一个操作：

1）标记该邀请码已被使用

2）把该邀请码跟注册账号进行关联

针对上述逻辑，首先针对库表进行设计和改造

直接在用户表新增了一个字段invitecode，目的就是把邀请码和注册者进行关联。

```go
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Mobile    string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
	InviteCode string      `json:"invitecode" gorm:"not null"`
}
```

然后，新建了一个表，专门用来存我们生成的邀请码，以及记录邀请码和生成邀请码的用户进行关联

```go
type InviteCode struct {
	gorm.Model
	Userid  uint      `json:"user_id" gorm:"not null"`
	Code    string      `json:"code" gorm:"not null"`
	Status  uint      `json:"status" gorm:"not null"`
}
```

库表建好了，下面就是在注册时，我们针对以上逻辑的代码实现：（全部代码请查看github）

```go

	if isTelephoneExist(DB,recommender){//查询这个推荐者存不存在
		id := getRecommenderId(DB,recommender)//先获取推荐者的ID
		//判断其邀请码的准确性
		var inviteCodes []model.InviteCode
		DB.Where("userid = ?",id).Find(&inviteCodes)//根据id去查这个推荐者有多少邀请码
		isValid := false
		for _,inviteCode := range inviteCodes{//轮寻他所有邀请码
			isValid = strings.EqualFold(code,inviteCode.Code)//一一对比，查看注册者输入的是否和库里的一致
			if  isValid {
				if inviteCode.Status == 0 {//如果一致，则再验证这个是否已经被用了
					DB.Model(&inviteCode).Update("status",1)
					break
				}else{
					response.Response(ctx,http.StatusOK,60002,nil,"邀请码已使用")
					return
				}
			}
		}
		if !isValid {
			response.Response(ctx,http.StatusOK,60002,nil,"邀请码不存在")
			return
		}
	}else {
    //最后这个提示，我是想避免被人拿这个接口测试用户存在与否
		response.Response(ctx,http.StatusOK,60002,nil,"邀请码不存在")
		return
	}
```

好了，到这里我们基本就完成这个简单的小功能了，直接go run main.go 去亲自体验一下吧。

ok just this.see you next...

