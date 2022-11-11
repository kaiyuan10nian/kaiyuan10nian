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
package FeatureTimeTools

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"strconv"
	"time"
)

/*
*时间戳和日期的相互转换功能实现
*params:
*app:由主界面传过来的fyne.App实例
 */
func TimeStampChange(app fyne.App) fyne.Window {
	myWindow := app.NewWindow("时间戳转换")
	title := canvas.NewText("时间戳转换", color.Black)
	content := container.New(layout.NewHBoxLayout(), title, layout.NewSpacer())
	timeNow := canvas.NewText("现在：", color.Black)
	clockTimeTamp := widget.NewLabel("")
	inputTimeStamp := widget.NewEntry()
	btnCopy := widget.NewButton("使用", func() {
		log.Println("tapped")
		inputTimeStamp.SetText(clockTimeTamp.Text)
	})

	updateTime(clockTimeTamp)
	row1 := container.New(layout.NewHBoxLayout(), timeNow, clockTimeTamp, btnCopy, layout.NewSpacer())
	timeStamp := canvas.NewText("时间戳：", color.Black)

	inputTimeStamp.SetPlaceHolder("输入时间戳")
	row2_1 := container.New(layout.NewGridWrapLayout(fyne.NewSize(120, 40)), inputTimeStamp)
	combo := widget.NewSelect([]string{"秒(s)", "毫秒(ms)"}, func(value string) {
		log.Println("Select set to", value)
	})
	combo.SetSelected("秒(s)")
	row2_2 := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 40)), combo)
	turnResult := widget.NewEntry()
	turnResult.SetPlaceHolder("转换结果")
	row2_3 := container.New(layout.NewGridWrapLayout(fyne.NewSize(160, 40)), turnResult)
	btnTurn := widget.NewButton("转换>", func() {
		var timetamp = ""
		if combo.Selected == "秒(s)" {
			timetamp = timeStamps2Date(inputTimeStamp.Text)
		} else {
			timetamp = timeStampMs2Date(inputTimeStamp.Text)
		}
		turnResult.SetText(timetamp)
	})
	beijingTimeStr := canvas.NewText("北京时间", color.Gray{0x99})
	row2 := container.New(layout.NewHBoxLayout(), timeStamp, row2_1, row2_2, btnTurn, row2_3, beijingTimeStr)

	timeDate := canvas.NewText("时间：", color.Black)
	inputDate := widget.NewEntry()
	inputDate.SetPlaceHolder("请输入日期")
	row3_1 := container.New(layout.NewGridWrapLayout(fyne.NewSize(160, 40)), inputDate)
	beijingTimeStr1 := canvas.NewText("北京时间", color.Gray{0x99})
	turnResult_date := widget.NewEntry()
	combo_date := widget.NewSelect([]string{"秒(s)", "毫秒(ms)"}, func(value string) {
		log.Println("Select set to", value)
	})
	combo_date.SetSelected("秒(s)")
	row3_2 := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 40)), combo_date)
	btnTurn_date := widget.NewButton("转换>", func() {
		var timetamp = ""
		if combo_date.Selected == "秒(s)" {
			timetamp = date2timeStamps(inputDate.Text)
		} else {
			timetamp = date2timeStampMs(inputDate.Text)
		}
		turnResult_date.SetText(timetamp)
	})

	turnResult_date.SetPlaceHolder("转换结果")
	row3_3 := container.New(layout.NewGridWrapLayout(fyne.NewSize(160, 40)), turnResult_date)
	row3 := container.New(layout.NewHBoxLayout(), timeDate, row3_1, beijingTimeStr1, btnTurn_date, row3_3, row3_2)

	myWindow.SetContent(container.New(layout.NewVBoxLayout(), content, row1, row2, row3))
	myWindow.Resize(fyne.NewSize(600, 300))
	go func() {
		for range time.Tick(time.Second) {
			updateTime(clockTimeTamp)
		}
	}()
	return myWindow
}

//更新时间戳
func updateTime(clock *widget.Label) {
	formatted := time.Now().Unix()
	timestamp := strconv.Itoa(int(formatted))
	clock.SetText(timestamp)
}

//秒时间戳转日期
func timeStamps2Date(timestamp string) string {
	timetamp, _ := strconv.Atoi(timestamp)
	timeTemplate := "2006-01-02 15:04:05"
	tm := time.Unix(int64(timetamp), 0)
	timeStr := tm.Format(timeTemplate)
	return timeStr
}

//日期转秒时间戳
func date2timeStamps(dateStr string) string {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //获取北京时间时区，很重要
	times, _ := time.ParseInLocation("2006-01-02 15:04:05", dateStr, TimeLocation)
	timeUnix := times.Unix()
	return strconv.Itoa(int(timeUnix))
}

//毫秒时间戳转日期
func timeStampMs2Date(timestamp string) string {
	timetamp, _ := strconv.Atoi(timestamp)
	timeTemplate := "2006-01-02 15:04:05"
	tm := time.UnixMilli(int64(timetamp))
	timeStr := tm.Format(timeTemplate)
	return timeStr
}

//日期转毫秒时间戳
func date2timeStampMs(dateStr string) string {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //获取北京时间时区，很重要
	times, _ := time.ParseInLocation("2006-01-02 15:04:05", dateStr, TimeLocation)
	timeUnix := times.UnixMilli()
	return strconv.Itoa(int(timeUnix))
}
