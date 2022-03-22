package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"kaiyuan10nian/common"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	InitConfig()//初始化配置
	db := common.InitDB()
	defer db.Close()
	InitGin()//初始化Gin框架并启动
}

func InitGin() {
	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")//这里加载配置文件中的端口
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(  workDir+ "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("")
	}
}
