package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"kaiyuan10nian/common"
	"kaiyuan10nian/dto"
	"kaiyuan10nian/model"
	"kaiyuan10nian/response"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func Register (ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	recommender := ctx.PostForm("recommender")
	code := ctx.PostForm("code")
	//数据验证
	if len(mobile) != 11 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		return
	}
	if len(name) == 0 {
		name = RandomString(10)
	}

	//判断手机号是否存在
	if isTelephoneExist(DB,mobile){
		response.Response(ctx,http.StatusOK,422,nil,"用户已经存在")
		return
	}
	//判断推荐人是否存在
	if isTelephoneExist(DB,recommender){
		id := getRecommenderId(DB,recommender)
		//判断其邀请码的准确性
		var inviteCodes []model.InviteCode
		DB.Where("userid = ?",id).Find(&inviteCodes)
		isValid := false
		for _,inviteCode := range inviteCodes{
			isValid = strings.EqualFold(code,inviteCode.Code)
			if  isValid {
				if inviteCode.Status == 0 {
					DB.Model(&inviteCode).Update("status",1)
					break
				}else{
					response.Response(ctx,http.StatusOK,60002,nil,"邀请码已使用")
					return
				}
			}
		}
		if !isValid {
			response.Response(ctx,http.StatusOK,60002,nil,"邀请码不存在1")
			return
		}
	}else {
		response.Response(ctx,http.StatusOK,60002,nil,"邀请码不存在2")
		return
	}


	//创建用户
	hasedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx,http.StatusInternalServerError,500,nil,"加密错误")
		return
	}
	newUser := model.User{
		Name: name,
		Mobile: mobile,
		Password: string(hasedPassword),
		InviteCode: code,
	}
	if err := DB.Create(&newUser).Error;err != nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,err.Error())
		return
	}
	//返回结果

	response.Success(ctx,nil,"注册成功")
}
//登录
func Login(ctx *gin.Context){
	ctx.Header("Access-Control-Allow-Origin", "*");
	DB := common.GetDB()
	//获取参数
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	//数据验证
	if len(mobile) != 11 {
		response.Response(ctx,http.StatusOK,422,nil,"手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx,http.StatusOK,422,nil,"密码不能少于6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	DB.Where("mobile = ?",mobile).First(&user)
	if user.ID == 0 {
		response.Response(ctx,http.StatusOK,422,nil,"用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password));err != nil{
		response.Response(ctx,http.StatusOK,400,nil,"密码错误")
		return
	}
	//发放token
	token,err := common.ReleaseToken(user)
	if err != nil{
		response.Response(ctx,http.StatusOK,500,nil,"系统异常")
		log.Printf("token generate error : %v",err)
		return
	}
	//返回结果
	response.Success(ctx,gin.H{"token":token},"登录成功")
}

func Info(ctx *gin.Context)  {
	user ,_:= ctx.Get("user")
	response.Success(ctx,gin.H{"user":dto.ToUserDto(user.(model.User))},"请求成功")
}
//创建用户邀请码
func Invite(ctx *gin.Context){
	udata ,_:= ctx.Get("user")
	dto := dto.ToUserDto(udata.(model.User))
	DB := common.GetDB()
	var inviteCodes []model.InviteCode
	DB.Where("userid = ?",dto.ID).Find(&inviteCodes)
	if len(inviteCodes) < 5{
		code := RandomString(5)
		newCode := model.InviteCode{
			Userid: dto.ID,
			Code: code,
			Status: 0,
		}
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

//判断手机号是否存在
func isTelephoneExist(db *gorm.DB, mobile string) bool {
	var user model.User
	db.Where("mobile = ?",mobile).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
//判断推荐人是否存在并获取其ID
func getRecommenderId(db *gorm.DB, mobile string) uint {
	var user model.User
	db.Where("mobile = ?",mobile).First(&user)
	return user.ID
}
//生成随机10个字符
func RandomString(n int) string {
	var letters = []byte("23456789qwertyupkjhgfdsazxcvbnmMNBVCXZASDFGHJKPOUYTREWQ")
	result := make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}