## <a id="9">项目部署</a>

我们开发的项目最终都会部署在服务器上去运行供所有人去访问，但是具体怎么部署呢？一直以来很多课程都是只讲项目怎么开发，极少人去讲怎么部署，因为部署又是一个很大的课题。很多公司把项目部署是放在运维下面的，但是我们作为一名开发人员不能被资本家腐蚀呀。不然你永远都是一颗螺丝钉，但凡换个地方你就蒙了。所以我们的目标是全栈，精通不精通的不重要，重要的是要懂每一个步骤中的道道。

那么这一节我们主要学习的就是怎么把项目部署到服务端让他跑起来，本节知识点比较多且杂。

#### 1、安装CentOS8 

准备工作：

​	1）准备一台电脑（可以把自己以前淘汰的电脑拿来用）

​	2）准备一个U盘

​	3）下载UItraISO(试用版即可)

第一步：使用UItraISO制作安装盘，这里需要注意一下，最好选择HDD++模式去写入，不然有可能在安装过程中找不到系统。

第二步：把制作好的U盘插电脑上进行系统安装，这里若是遇到安装失败的情况，请参考以下解决方案：

```go
选中第一项后按键盘上的Tab键进行编辑（以前版本按E键进入编辑界面）。
重点来了：将vmlinuzinitrd=initrd.img inst.stage2=hd:LABEL=CentOS\x207\x20x86_64 quiet 

改为vmlinuz initrd=initrd.img inst.stage2=hd:LABEL=U盘名称 quiet

或改为vmlinuz initrd=initrd.img inst.stage2=hd:/dev/sdc4 quiet

Sdc4：使用 linux dd 查看（可不是所有设备的都是sdc4,请直接cd /dev然后ls查看你的U盘属于哪一个）

然后按Esc退出编辑执行安装或Ctrl+x执行安装（不同版本略有不同）。
```

第三步：能来到第三步就非常厉害了，给你点赞。在这一步主要就是设置安装目标，其他什么时间、软件选择、设置密码啥的都很easy。这里会遇到一个坑，就是在设置分区的时候你会发现不知道怎么选，要么空闲不够要么不知道点哪里。请看下图，框框中有几个你就选择几个，全选中后下面选择自定义，然后直接点左上角的Done就会直接进入分区界面了。

![WechatIMG289.png](https://s2.loli.net/2022/03/30/lIALW3BveftyrR5.png)

具体分区请按照你自己的硬盘大小自己分配，但是以下几个必须配置且类型必须选对，不然后面会直接影响你安装是否成功。

| 文件名称 | 文件系统类型 | 存储类型 |
| -------- | ------------ | -------- |
| /boot    | ext4         | 标准分区 |
| /        | Xfs          | LVM      |
| /home    | Xfs          | LVM      |
| /tmp     | Xfs          | LVM      |
| Swap     | Xfs          | 标准     |

一切搞定后确定，然后弹窗提示硬盘原来数据会被抹除，直接确定即可。

回到首页后，点右下角的安装即可。大概等个5分钟吧，安装完成。

#### 2、CentOS源配置

CentOS 8操作系统版本结束了生命周期（EOL），Linux社区已不再维护该操作系统版本。所以如果你不改的话你会发现什么都是安装失败，直接参考以下配置操作即可：(以下操作只针对CentOS8，其他版本请绕行)

```sh
1. 备份
mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
2. 下载新的 CentOS-Base.repo 到 /etc/yum.repos.d/
centos8（centos8官方源已下线，建议切换centos-vault源）
wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-vault-8.5.2111.repo
或者
curl -o /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-vault-8
```

修改后，yum命令就可以使用了。

1、首先设置自己本地语言系统，不然每次都是几个白色小方块你根本就不知道是什么错误信息

```shell
下载的centos8镜像系统默认使用的是ISO/IEC 15897字符集
需要改成UTF-8.
先安装所有的字符集
dnf install langpacks-en glibc-all-langpacks -y
设置字符集
localectl set-locale LANG=en_US.UTF-8

```

2、由于我们安装的系统没有net-tools包所以用不了ifconfig等命令，那肯定不能忍：

```shell
yum install net-tools
```

3、后面我们肯定会用到数据库吧，我喜欢用MySQL，直接安装：

```
sudo dnf install @mysql
```

4、我们后端用的golang，那顺便也安装一下吧；

```shell
yum install go -y
```

#### 3、编译要部署的Go项目

因为是demo，所以功能非常简单，整个项目只有一个main.go，其内容如下：

```go
func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "first example.")
	})

	router.Run()
}
```

先试一下本地是否能跑起来，母庸质疑肯定可以了。

下面直接打包：

```go
env GOOS=linux GOARCH=386 go build main.go
```

这里需要说一下，因为我的开发环境是mac,而我要部署的目标服务器是linux的，所以GOOS这里我选择了linux,不然你编译后会发现根本运行不起来。另外这里的386就是系统环境标识，具体有点复杂后面有机会再讨论。

命令执行后3S吧，你就会在你的项目根目录下看到一个main文件，直接把这个文件传到目标服务器中。

#### 4、运行你的go项目

当你把main文件传到你的服务器后，你直接在命令行输入   ./main    即可。

但是一般这个时候会提示你没有权限，直接使用：chmod 777 main，这样你就有权限了，继续输入上面的命令。

运行后，你会看到一句提示语： Listening and serving HTTP on :8080.这个时候就表示你运行成功了。

#### 5、关闭防火墙

一般这个时候你肯定会迫不及待的去另一台电脑上打开浏览器输入：ip:8080去享受此刻的胜利果实。但不幸的是你发现你访问不了，这是因为centos8自带防火墙且默认是开启的，这个时候只需要去关闭即可。下面是我收集的常用命令请收藏：

```shell
#进程与状态相关
 2 systemctl start firewalld.service            #启动防火墙  
 3 systemctl stop firewalld.service             #停止防火墙  
 4 systemctl status firewalld                   #查看防火墙状态
 5 systemctl enable firewalld             #设置防火墙随系统启动
 6 systemctl disable firewalld                #禁止防火墙随系统启动
 7 firewall-cmd --state                         #查看防火墙状态  
 8 firewall-cmd --reload                        #更新防火墙规则   
 9 firewall-cmd --list-ports                    #查看所有打开的端口  
10 firewall-cmd --list-services                 #查看所有允许的服务  
11 firewall-cmd --get-services                  #获取所有支持的服务  
12 
13 #区域相关
14 firewall-cmd --list-all-zones                    #查看所有区域信息  
15 firewall-cmd --get-active-zones                  #查看活动区域信息  
16 firewall-cmd --set-default-zone=public           #设置public为默认区域  
17 firewall-cmd --get-default-zone                  #查看默认区域信息  
18 
19 
20 #接口相关
21 firewall-cmd --zone=public --add-interface=eth0  #将接口eth0加入区域public
22 firewall-cmd --zone=public --remove-interface=eth0       #从区域public中删除接口eth0  
23 firewall-cmd --zone=default --change-interface=eth0      #修改接口eth0所属区域为default  
24 firewall-cmd --get-zone-of-interface=eth0                #查看接口eth0所属区域  
25 
26 #端口控制
27 firewall-cmd --query-port=8080/tcp             # 查询端口是否开放
28 firewall-cmd --add-port=8080/tcp --permanent               #永久添加8080端口例外(全局)
29 firewall-cmd --remove-port=8800/tcp --permanent            #永久删除8080端口例外(全局)
30 firewall-cmd --add-port=65001-65010/tcp --permanent      #永久增加65001-65010例外(全局)  
31 firewall-cmd  --zone=public --add-port=8080/tcp --permanent            #永久添加8080端口例外(区域public)
32 firewall-cmd  --zone=public --remove-port=8080/tcp --permanent         #永久删除8080端口例外(区域public)
33 firewall-cmd  --zone=public --add-port=65001-65010/tcp --permanent   #永久增加65001-65010例外(区域public) 

```

#### 6、后台运行配置

经过以上步骤，其实你已经可以在浏览器对这个服务进行访问了。可这个时候你又发现一个问题，在你关闭当前与服务器的连接后，这个服务又访问不了了。不用担心，请看下面：

```shell
CentOS后台运行和关闭、查看后台任务命令
一、&
加在一个命令的最后，可以把这个命令放到后台执行，如
watch -n 10 sh test.sh & #每10s在后台执行一次test.sh脚本
二、ctrl + z
可以将一个正在前台执行的命令放到后台，并且处于暂停状态。
三、jobs
查看当前有多少在后台运行的命令
jobs -l选项可显示所有任务的PID，jobs的状态可以是running, stopped, Terminated。但是如果任务被终止了（kill），shell 从当前的shell环境已知的列表中删除任务的进程标识。
四、fg
将后台中的命令调至前台继续运行。如果后台中有多个命令，可以用fg %jobnumber（是命令编号，不是进程号）将选中的命令调出。
五、bg
将一个在后台暂停的命令，变成在后台继续执行。如果后台中有多个命令，可以用bg %jobnumber将选中的命令调出。
六、kill
法子1：通过jobs命令查看job号（假设为num），然后执行kill %num
法子2：通过ps命令查看job的进程号（PID，假设为pid），然后执行kill pid
前台进程的终止：Ctrl+c
七、nohup
如果让程序始终在后台执行，即使关闭当前的终端也执行（之前的&做不到），这时候需要nohup。该命令可以在你退出帐户/关闭终端之后继续运行相应的进程。关闭中断后，在另一个终端jobs已经无法看到后台跑得程序了，此时利用ps（进程查看命令）
ps -aux | grep “test.sh” #a:显示所有程序 u:以用户为主的格式来显示 x:显示所有程序，不以终端机来区分

```

至此，开始享受最后的胜利果实吧~

OK,Just it,See you next...

