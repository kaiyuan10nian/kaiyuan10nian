## <a id="4">Go语言下的Gin详解及Demo实践</a>

今天没有去继续更新开源10年的项目，因为昨天接触到几个新的知识点，所以今天加强一下对他们的认识，下面是本节的一个知识点。

1）Gin Web框架的认识

2）Gin加载静态资源

3）Gin加载动态资源

在学习Gin的过程中动手搞了一个小Demo，把Gin的简单用发都跑了一下，强烈建议各位读者不要只看文章，自己动手写一下效果更佳。

项目在GitHub上的地址：[https://github.com/kaiyuan10nian/GinDemo](https://link.zhihu.com/?target=https%3A//github.com/kaiyuan10nian/GinDemo)



下面开始我们今天的知识点分享。

一、创建一个Go web项目，这里我命名为：GinDemo，方便我们的学习。

项目结构

![img](https://pic4.zhimg.com/80/v2-e348197212d67e1a12cb7e024000d053_720w.jpg)

Gin安装：go get -u [http://github.com/gin-gonic/gin](https://link.zhihu.com/?target=http%3A//github.com/gin-gonic/gin)

直接在Goland下面的Terminal中输入就可以了。看到下图就表示你安装成功了

![img](https://pic4.zhimg.com/80/v2-1424c34bf85b347ed21d21e8f93a572b_720w.jpg)

二、先写一个小案例，并在Postman或者浏览器中打开看下效果。

在GinDemo下面创建一个main.go文件，输入以下内容：

![img](https://pic2.zhimg.com/80/v2-da8f3583ed46a2b7301de6cd7279c471_720w.jpg)

然后在下面的Termonal中输入以下指令运行该项目：

Go run main.go

看到下图即表示你已经运行起来了。

![img](https://pic3.zhimg.com/80/v2-8ce43416c230aa18159f8a8059cfc3b6_720w.jpg)

这个时候，打开你的浏览器，输入：localhost:8080/ping,将显示以下内容：

![img](https://pic4.zhimg.com/80/v2-faefae54729a836ba403d187b8ff2c7b_720w.jpg)

整个Demo中的注释还是比较清楚的，每一行代码是什么意思，有什么作用等都比较简单，运行到这里基本上Gin的精髓就已经学到了。

下面两个知识点是可以解决我们在真正的实际项目中经常会遇到的需求的。所以我这里单独拉了出来写一下。当然，它还有其它别的更多用法，我们就不一一说了，至少掌握了这两个对于普通的开发工作就足以应对了。

三、加载静态资源

在实际项目中我们经常会用到很多静态资源，比如：图片、文件等。那Gin是怎么处理的呢？还是看案例，下面我们对上面的main.go进行一下简单修改：

![img](https://pic2.zhimg.com/80/v2-33dc491bd60886def79d33e3d2b6e1e1_720w.jpg)

主要加了两行代码，用到两个函数。注释中已经描述的很清楚了，说一下代码中未描述的内容，非常重要，这两个函数的第一个参数就是相对地址，也就是说是用户端访问的时候访问的地址，第二个参数是本地服务器的地址，也就是引用的地址。

比如我们运行上面项目后，若想访问/Users/fu/GolandProjects/GinDemo/web/static里面的图片，那么直接如下操作即可：

![img](https://pic3.zhimg.com/80/v2-6ddc155caafbdca35df8f7bd6462938e_720w.jpg)

若想加载第二种加载形式的，则直接这么访问：

![img](https://pic1.zhimg.com/80/v2-ee6f5f357df3e6e5828b2276d705db3c_720w.jpg)

其实，这两种访问形式访问的都是同一张图片。

这就是静态加载资源的方式方法，那么我们做一个web站点肯定不是仅仅有静态资源，肯定还有动态资源，动态资源的加载怎么实现？

四、加载动态资源

首先，在GinDemo-web下新建templete并在其中新增一个index.html文件，文件内容很简单，这里制作演示所以就没有去做接口互动。

![img](https://pic4.zhimg.com/80/v2-8c450de1340e43c72fc1c9cf84dfd963_720w.jpg)

然后，在根目录下创建了controller文件夹并创建con.go，主要用来存放逻辑层的操作，受JAVA开发的影响我这么去做了，你可以随意哈，怎么高兴怎么来。

![img](https://pic2.zhimg.com/80/v2-69c4988cfa34759d472d8a883f0f7b79_720w.jpg)

还是在main.go中做修改,并加载trmplete下面的资源，然后做了一个web分组，分组中仅有一个接口，其处理放到了con.go文件中的IndexController函数中。

![img](https://pic3.zhimg.com/80/v2-fb6809244ff085d3c8cb2ca599f35fda_720w.jpg)

直接运行，然后浏览器中访问index.html看看有什么效果？

![img](https://pic1.zhimg.com/80/v2-93e4ebd52e4940966c9b330d2e3f0570_720w.jpg)

到这里，其实这个项目就算完成了，其中涉及到的知识点我们也都了解了，就这些内容已经完全够现阶段的我们使用了。

拜拜。。。See you next.







---









