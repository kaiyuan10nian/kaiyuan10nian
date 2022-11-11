/*
Copyright (c) [2022] [开源十年]
[开源十年] is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"kaiyuan10nian/ToolsBaBa/FeatureTimeTools"
	"os"
)

/*
工具集主界面
*/
func main() {
	os.Setenv("FYNE_FONT", "./ToolsBaBa/res/msyhl.ttc") //设置env环境
	myApp := app.New()
	mainWindow := myApp.NewWindow("开源十年")
	//挂载时间戳转换功能
	btnTimeStamp := widget.NewButton("时间戳转换", func() {
		FeatureTimeTools.TimeStampChange(myApp).Show()
	})
	//挂载BASE64加解密功能
	//Todo
	//挂载AES加解密功能
	//Todo
	//挂载DES加解密功能
	//Todo
	//挂载JSON转换功能
	//Todo
	//挂载RSA加解密功能
	//Todo
	//挂载Unicode功能
	//Todo
	content := container.New(layout.NewHBoxLayout(), btnTimeStamp, layout.NewSpacer())
	tabName := canvas.NewText("工具集", color.Black)
	tabLine := canvas.NewLine(color.Gray{0x99})
	mainWindow.SetContent(container.New(layout.NewVBoxLayout(), tabName, tabLine, content))
	mainWindow.Resize(fyne.NewSize(600, 300))
	mainWindow.ShowAndRun()
	os.Unsetenv("FYNE_FONT")
}
