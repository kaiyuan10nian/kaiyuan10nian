## <a id="3">Blog项目骨架搭建---Go语言（第一阶段）</a>

今天周六，利用两个多小时的时间把昨天的需求简单梳理了一下，并开始搭建项目，进行了简单的编辑，到目前为止已经实现注册功能了。

这个项目到今天是第三天了，我并不能保证每天都更新日记，但确实每天都会更新项目，时间多久多做点，时间少就少做点，总之这个事情是不可以断的。希望各位也给我助助威，鞭策一下我这颗懒惰的心。

本节主要涉及以下几个知识点：

1）gin web框架

2）viper 配置管理

3）Gorm 数据库连接

这三个知识点我后面单独拉出来介绍，这节仅简单介绍并陈述我都做了什么。

**一、选择Gin框架**

在 Go语言开发的 Web 框架中，有两款著名 Web 框架分别是 Martini 和 Gin，两款 Web 框架相比较的话，Gin 自己说它比 Martini 要强很多。

Gin 是 Go语言写的一个 web 框架，它具有运行速度快，分组的路由器，良好的崩溃捕获和错误处理，非常好的支持中间件和 json。总之在 Go语言开发领域是一款值得好好研究的 Web 框架

gin安装：go get -u [http://github.com/gin-gonic/gin](https://link.zhihu.com/?target=http%3A//github.com/gin-gonic/gin)   //错误示例 多谢@changpingdengdeng的指正

gin安装：go get -u github.com/gin-gonic/gin  //正确示例

本项目中初始化位置：

![img](https://pic4.zhimg.com/80/v2-55427a3d55a310224ef9b225e8e19b33_720w.jpg)

**二、选自viper做配置管理**

viper 是一个配置解决方案，拥有丰富的特性：

- 支持 JSON/TOML/YAML/HCL/envfile/Java properties 等多种格式的配置文件；
- 可以设置监听配置文件的修改，修改时自动加载新的配置；
- 从环境变量、命令行选项和io.Reader中读取配置；
- 从远程配置系统中读取和监听修改，如 etcd/Consul；
- 代码逻辑中显示设置键值。

Viper安装：go get [http://github.com/spf13/viper](https://link.zhihu.com/?target=http%3A//github.com/spf13/viper)

本项目中初始化位置：

![img](https://pic1.zhimg.com/80/v2-69a3df27012e7d4a8c9d86fe13462e4c_720w.jpg)

![img](https://pic4.zhimg.com/80/v2-37fdd5d36f58bd4bd6834e5d55718e5b_720w.jpg)

根据以上代码可知，我的配置文件放下项目路径/config下，文件为：application.yml。

**三、数据库连接使用Gorm库**

gorm是go语言的一个orm框架,具体的原理及思想我也介绍不清楚，你只需要知道它是你操作数据库的桥梁即可、

Gorm安装：go get -u [http://github.com/jinzhu/gorm](https://link.zhihu.com/?target=http%3A//github.com/jinzhu/gorm)

本项目中初始化位置：

![img](https://pic2.zhimg.com/80/v2-9d5fcd6309a98e28275899295092fd5d_720w.jpg)

![img](https://pic2.zhimg.com/80/v2-7b48b6ed55166ffcfe9c3f16dd1105ed_720w.jpg)

上图中的viper.GetString("datasource.driverName")就是利用前面说过的配置管理viper去配置文件中获取相对应的参数。

具体的配置文件如下：

![img](https://pic1.zhimg.com/80/v2-86b3083dfe913f9c73df4547c613a7c4_720w.jpg)

这里留一个小彩蛋，只有在真正运行项目的时候才会发现哦

**四、创建数据模型**

一共创建了两个数据模型

![img](https://pic2.zhimg.com/80/v2-d08876977c311cd726cefe9d24d12225_720w.jpg)

![img](https://pic2.zhimg.com/80/v2-e4749fcf00652bd4579e6b38213a3fb9_720w.jpg)

这里一共设计到三个知识点

1）gorm.Model这个标记一个结构中有一个结构中没有，那他是用来干什么的呢？其实这就相当于是继承，加上这个后就相当于继承了Model，不加这个代表不继承Model。而Model结构中有以下四个定义好的字段：

![img](https://pic2.zhimg.com/80/v2-5a22f68ffcd8adf236cffa6de1e82ce9_720w.jpg)

这是Gorm自带的，所以你清楚这一点就行了。

2）上面Article结构体中有一个uuid.UUID标记，这个是哪里来的？其实是引用自：

uuid"[http://github.com/satori/go.uuid](https://link.zhihu.com/?target=http%3A//github.com/satori/go.uuid)"

它的主要功能就是在实际项目中，经常会使用到一个唯一标识的，比如唯一标识一条记录等情况，这个go.uuid项目库就是干这个事情的。

3）在上面Article结构体中还有一个非自带类型：Time，这个是自定义的，其作用就是把时间格式化了而已。如果不格式化的话，它显示的就是时间戳，这个大家应该都知道吧？所以它就是让我们能更方便的查看时间。

**五、写了一个接口**

账号注册 v1/account/register

![img](https://pic1.zhimg.com/80/v2-ce76f08369772220496c42bbdc9b59a4_720w.jpg)

这节就先开这一个接口吧，点到为止。贪多贪快嚼不烂，我们的目的是掌握开发过程中出现的各个知识点，并不是完成功能的开发。所以大家不要慌尽量把这节的内容摸透咽下去再进行下一步的开发。

下一节安排：

1、Go 项目实战 之 Gin框架的详解

2、Go 项目实战 之 配件管理viper 详解

3、Go 项目实战 之 数据库连接Grom详解

这个项目到目前为止，基本上骨架就出来了，当然项目本身问题还是很多的，我们会在后面一步一步去完善，为的就是在完善中学习。直接一步到位的框架设计只会在外包项目中出现，我们又不赶进度，慢慢来哈。



