哎呀呀，这次更新间隔拉得有点久，是在是对不住各位支持我的小伙伴，最近公司刚启动了一个新项目，时间实在是紧张，希望各位多多见谅。下面开始今天的分享正文。

**Go项目实战之日志必备篇**

你在某公司技术部经常听到的一句话就是：“稍等，让我查下**日志**再回复你。”

日志---其实就是项目在运行期间留下的痕迹。

就好比你冬天去打兔子，第一件事就是找雪地上兔子留下的脚印，然后顺着脚印去找到兔子的老窝，最后满载而归。

项目开发也是一样的道理，你要在项目开发期间想尽办法去让业务逻辑在运行期间能留下更多的关键信息，这样在项目正常运行后就会留下它运行的痕迹，你就可以通过这个痕迹快速寻找到问题的根源，从而一网打尽。

而在Go语言中只提供了一个标准库log.

不需要安装即可使用，它是一个非常小巧的日志库，大家有空可以去看看。

log只提供了三个简单的接口，对于某些大型项目来说有点太鸡肋，所以本篇我要介绍的不是log而是它的哥哥logrus。

logrus是一个完全兼容log的标准库，目前在GitHub上已经超20K stars了。

它支持两种日志输出格式：文本 And JSON

这对于想通过日志去做更多数据分析的项目来说简直是太爽了。

本篇分两部分：

- 1、logrus的基本用法介绍

- 2、封装logrus用于**开源十年**项目中

​	  	2.1)关于使用lumberjack对logrus生成的日志进行分包

1、logrus的基本用法介绍

1）安装

```go
go get -u github.com/sirupsen/logrus
```

2）设置日志输出等级

​	老程序员都知道，一个项目完整的开发周期分好几个阶段，所以在不同的阶段我们需要的日志信息也是不一样的。所以在使用logrus的时候我们要去设置其日志输出的等级,从而筛选出我们需要的信息内容。

​	那么在设置日志输出等级之前，我们需要了解logrus共区分了几个等级？

| 级别 | 等级              | 解释                                                 |
| ---- | ----------------- | ---------------------------------------------------- |
| 1    | logrus.TraceLevel | 非常小粒度的信息                                     |
| 2    | logrus.DebusLevel | 一般程序中输出的调试信息                             |
| 3    | logrus.InfoLevel  | 关键操作（核心流程日志）                             |
| 4    | logrus.WarnLevel  | 警告信息                                             |
| 5    | logrus.ErrorLevel | 错误信息                                             |
| 6    | logrus.FatalLevel | 致命错误，出现后程序无法运行，输出日之后程序停止运行 |
| 7    | logrus.PanicLevel | 记录日志，然后panic                                  |

左边有个**级别** 大家一定要记清楚顺序，因为在logrus中，高于设置级别的日志是不会输出的，默认设置级别是InfoLevel

示例：

```go
logrus.SetLevel(logrus.TraceLevel)
logrus.Trace("1---trace---msg")
logrus.Debug("2---debug---msg")
logrus.Info("3---info---msg")
logrus.Warn("4---warn---msg")
logrus.Error("5---error---msg")
logrus.Fatal("6---fatal---msg")
logrus.Panic("7---panic---msg")
```

运行之后我们看下日志输出情况：

```go
TRAC[0000] 1---trace---msg                              
DEBU[0000] 2---debug---msg                              
INFO[0000] 3---info---msg                               
WARN[0000] 4---warn---msg                               
ERRO[0000] 5---error---msg                              
FATA[0000] 6---fatal---msg  
```

如果上面代码中我们把：

```go
logrus.SetLevel(logrus.TraceLevel)
```

修改为：

```go
logrus.SetLevel(logrus.InfoLevel)
```

然后，再运行看下输出情况：

```go
INFO[0000] 3---info---msg                               
WARN[0000] 4---warn---msg                               
ERRO[0000] 5---error---msg                              
FATA[0000] 6---fatal---msg  
```

可以看到，比info级别低的就不再输出了。

3）在日志中输出具体文件和函数位置

为了快速定位问题根源，很多时候我们会在调试阶段直接把文件路径和函数直接输出出来，这样我们就不需要再去定位寻找了，解决问题的效率将会得到大大的提升。

logrus提供了专门的配置，只需要在初始化logrus的时候调用SetReportCaller()函数并设置为true即可。

示例：

```go
	logrus.SetReportCaller(true)
	logrus.Info("3---info---msg")
```

直接运行看效果：

```go
INFO[0000]/Users/fu/GolandProjects/logrusDemo/main.go:29 main.main() 3---info---msg    
```

4）添加附属信息

我们做为后端开发人员，时刻把并发的问题放在心头是本能。所以在记录日志时，你可能也会思考怎么去区分日志。

比如：

哪些日志是张三留下的？哪些日志是李四留下的？为什么同样的逻辑流程张三和李四输出的结果不一样呢？

这个时候你或许在想，如果我给这些日志打上“张三”“李四”的备注是不是就好找多了？

logrus提供了解决方案，就是WithField和WithFields ,允许在输出中添加一些字段，比如：

```go
logrus.WithFields(logrus.Fields{
		"UUID": "12345678",
	}).Info("info msg")
```

日志输出：

```go
INFO[0000] 3---info---msg                                UUID=12345678
```

这是针对单个的使用方式，如果想批量使用更好办：

```go
requestLogger := logrus.WithFields(logrus.Fields{
		"UUID": "12345678",
	})
requestLogger.Info("3---info---msg")
requestLogger.Error("5---error---msg")
```

日志输出：

```go
INFO[0000] 3---info---msg                                UUID=12345678
ERRO[0000] 5---error---msg                               UUID=12345678
```

5）JSON格式输出日志

上面我们输出日志的时候用的是默认的输出格式，也就是文本格式。但是在很多业务中我们做数据统计或者数据分析的时候依靠的源数据都是日志，如果是文本格式那么用起来就不是那么的顺手，换成json格式的话会不会带来很大的方便呢？

logrus不同于log的最大之处就是提供了json格式的输出，只需要在初始化的时候设置SetFormatter即可。

```go
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Trace("1---trace---msg")
	logrus.Debug("2---debug---msg")
	logrus.Info("3---info---msg")
	logrus.Warn("4---warn---msg")
	logrus.Error("5---error---msg")
	logrus.Fatal("6---fatal---msg")
	logrus.Panic("7---panic---msg")
```

跟1）中的实例一样，只是添加了logrus.SetFormatter(&logrus.JSONFormatter{})，下面我们看下输出的格式：

```go
{"level":"trace","msg":"1---trace---msg","time":"2022-05-14T11:37:56+08:00"}
{"level":"debug","msg":"2---debug---msg","time":"2022-05-14T11:37:56+08:00"}
{"level":"info","msg":"3---info---msg","time":"2022-05-14T11:37:56+08:00"}
{"level":"warning","msg":"4---warn---msg","time":"2022-05-14T11:37:56+08:00"}
{"level":"error","msg":"5---error---msg","time":"2022-05-14T11:37:56+08:00"}
{"level":"fatal","msg":"6---fatal---msg","time":"2022-05-14T11:37:56+08:00"}
```

ok,到这里logrus的基本操作我们就明白了，下面针对在开源十年项目中我们怎么去系统的运用它。

2、封装logrus用于**开源十年**项目中

我封装了一个logger.go的文件，放在了config目录下面，下面我把代码完整放上来，然后在备注中去一一解释一下。

```go
package config

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

var logger *logrus.Logger
//日志名称
const (
	//日志文件名
	LOG_NAME = "kaiyuanshinian"
	//日志文件后缀
	LOG_SUFFIX = ".log"
	//单个日志文件大小，单位MB
	LOG_SIZE = 50
	//日志文件个数
	LOG_BACKUP = 10
	//日志文件最大天数
	LOG_DATE = 7
)

//设置日志输出到文件
func setOutPut(log *logrus.Logger, log_file_path string) {
	logconf := &lumberjack.Logger{
		Filename:   log_file_path,
		MaxSize:    LOG_SIZE,   // 日志文件大小，单位是 MB
		MaxBackups: LOG_BACKUP, // 最大过期日志保留个数
		MaxAge:     LOG_DATE,   // 保留过期文件最大时间，单位 天
		Compress:   true,       // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
	}
	log.SetOutput(logconf)
}

//初始化日志模块
func InitLogger() {
	log_file_path := path.Join("./", LOG_NAME+LOG_SUFFIX)
	logger = logrus.New()
	setOutPut(logger, log_file_path)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
//获取logrus操作对象
func GetLogger() *logrus.Logger {
	return logger
}

//gin请求消息也写入日志
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()               // 开始时间
		c.Next()                              // 处理请求
		endTime := time.Now()                 // 结束时间
		latencyTime := endTime.Sub(startTime) // 执行时间
		reqMethod := c.Request.Method         // 请求方式
		reqUri := c.Request.RequestURI        // 请求路由
		statusCode := c.Writer.Status()       // 状态码
		clientIP := c.ClientIP()              // 请求IP
		logger.Infof("| %3d | %13v | %15s | %s | %s", statusCode, latencyTime, clientIP, reqMethod, reqUri ) // 日志格式
	}
}
```

上面是对logrus的封装，大家应该都看的明白的我就不一一啰嗦了，那么怎么使用呢？（上面代码中只有lumberjack是我们之前没有提及过得，下面解释）

1）初始化

直接在项目启动的时候把logrus的初始化加进去即可

```go
func InitConfig() {
	config.InitLogger()//初始化logrus
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(""+err.Error())
	}
}
```

2）使用

使用就更简单了，直接在项目需要的位置进行调用即可。

2.1）route中的使用

```go
func CollectRoute(r *gin.Engine) *gin.Engine  {
	r.Use(config.LoggerToFile())//添加日志记录
	r.POST("/v1/account/register", controller.Register)
	
	...//省略了一些代码，都是以前写的，项目中有

	return r
}
```

2.2）逻辑中的使用

```go
config.GetLogger().Debugf("aaaaa")
```

3）lumberjack

上面封装的代码中大家发现多了个新东西lumberjack，他是干啥用的呢？

对，切分日志文件的。

有的时候我们的日志需要大量的去记载，如果都记录在一个文件中，万一发生点什么不可描述的事情导致文件丢失了那我们就只有哭的份了。

所以一是为了安全，二是为了方便，我们要针对日志文件进行分割保存。

```go
logconf := &lumberjack.Logger{
		Filename:   log_file_path,
		MaxSize:    LOG_SIZE,   // 日志文件大小，单位是 MB
		MaxBackups: LOG_BACKUP, // 最大过期日志保留个数
		MaxAge:     LOG_DATE,   // 保留过期文件最大时间，单位 天
		Compress:   true,       // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
	}
	log.SetOutput(logconf)
```

它的使用非常简单，设置好参数，在SetoutPut中传入进去就可以了，当文件大于我们设定的MaxSize时，会自动进行分割保存。

ok ,just this...

今天就先写这么多吧。谢谢大家的支持哦~