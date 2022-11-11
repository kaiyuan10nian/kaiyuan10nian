package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"kaiyuan10nian/Blog/config"
	"kaiyuan10nian/Blog/response"
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
			response.Response(ctx, http.StatusOK, 50001, nil, "文件类型不合法")
			return
		}
		//3、创建图片保存目录,linux下需要设置权限（0755可读可写） kaiyuan/upload/image20220915
		currentTime := time.Now().Format("20060102")
		// 生成目录文件夹，并错误判断
		if err := os.MkdirAll("/opt/server/nginx-1.18/html/kaiyuan/upload/"+currentTime, 0755); err != nil {
			config.GetLogger().Error(err, "上传错误", false)
			// 返回值
			response.Response(ctx, http.StatusOK, 50001, nil, "MkdirAll失败")
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
			response.Response(ctx, http.StatusOK, 5001, nil, "文件保存失败")
			return
		}
		imageurl := strings.Replace(saveDir, "/opt/server/nginx-1.18/html", "https://xiaoyin.live", -1)
		// 返回值
		response.Success(ctx, gin.H{"imageurl": imageurl}, "上传成功")
		return
	}
}
