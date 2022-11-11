package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"kaiyuan10nian/Blog/common"
	"kaiyuan10nian/Blog/config"
	"kaiyuan10nian/Blog/route"
	"net/http"
)

func main() {
	InitConfig() //初始化配置
	db := common.InitDB()
	defer db.Close()
	InitGin() //初始化Gin框架并启动
}

func InitGin() {
	r := gin.Default()
	r = route.CollectRoute(r)
	r.StaticFS("/kaiyuan", http.Dir("/opt/server/nginx-1.18/html/kaiyuan"))
	port := viper.GetString("server.port") //这里加载配置文件中的端口
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}

func InitConfig() {
	config.InitLogger()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic("" + err.Error())
	}
}
