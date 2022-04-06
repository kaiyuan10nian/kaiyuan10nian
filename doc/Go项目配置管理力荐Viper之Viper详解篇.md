## <a id="5">Go项目配置力荐Viper之Viper详解篇</a>

引言：今天还是补充上次我们搭建项目是遇到的知识点Viper，以下内容了解后基本上在使用Viper上已经无压力了，后面进阶的时候我再补充Viper的进阶篇。

---



Viper 是一个完整的 Go 应用程序配置解决方案，优势就在于开发项目中你不必去操心配置文件的格式而是让你腾出手来专注于项目的开发。其特性如下：

- 支持 JSON/TOML/YAML/HCL/envfile/Java properties 等多种格式的配置文件；
- 可以设置监听配置文件的修改，修改时自动加载新的配置；
- 从环境变量、命令行选项和io.Reader中读取配置；
- 从远程配置系统中读取和监听修改，如 etcd/Consul；
- 代码逻辑中显示设置键值

**注：Viper让需要重启服务器才能使配置生效的日子一去不复返！！！**这才是VIper最大的魅力

#### 基础配置

Viper没有默认的基础配置，所以在使用的过程中我们初始化Viper实例的时候需要告诉Viper你的配置路径、配置格式、配置名称等等信息。Viper虽然支持多配置同时使用，但是一个Viper实例只能寻一个配置路径。

示例：

```Go
viper.SetConfigName("config") // 配置文件名 (不带扩展格式)
viper.SetConfigType("yaml") // 如果你的配置文件没有写扩展名，那么这里需要声明你的配置文件属于什么格式
viper.AddConfigPath("/etc/appname/")   // 配置文件的路径

err := viper.ReadInConfig() //找到并读取配置文件
if err != nil { // 捕获读取中遇到的error
	panic(fmt.Errorf("Fatal error config file: %w \n", err))
}
```

如果在项目中你的配置文件找不到或者找的过程中出error了，怎么办？可以参考下面这么做：

```go
if err := viper.ReadInConfig(); err != nil {
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// 配置文件没有找到，Todo
	} else {
		// 配置文件找到了，但是在这个过程有又出现别的什么error,Todo
	}
}

// 配置文件成功找到，并未再出现什么error则继续往下面执行
```

#### 写入运行时配置

很多时候我们需要记录程序在运行时的一些配置参数，那么Viper也是可以做到的，在Viper中提供了四个运行时记录配置的方法：

- WriteConfig：将当前Viper配置写入指定路径，如果保存路径不存在则报错，存在覆盖
- SafeWriteConfig：将当前Viper配置写入指定路径，如果保存路径不存在则报错，存在则不会覆盖
- WriteConfigAs：将当前Viper配置写入指定路径，覆盖指定文件（如果路径存在的话）
- SafeWriteConfigAs：将当前Viper配置写入指定路径，除指定文件外其他的都覆盖（如果路径存在的话）

```go
viper.WriteConfig() // 写入当前配置到'viper.AddConfigPath()' 和 'viper.SetConfigName'设定的路径
viper.SafeWriteConfig()
viper.WriteConfigAs("/path/to/my/.config")
viper.SafeWriteConfigAs("/path/to/my/.config") // 这里将报错，因为已经写入过了
viper.SafeWriteConfigAs("/path/to/my/.other_config")
```

#### 如何让配置实时生效

Viper支持应用程序在运行时实时读取配置文件的能力，只需要告诉Viper去watchConfig即可。或者直接自己去实现一个函数，每次更改配置文件时都运行。

示例：（该示例运行前请确保你已经设置了configPaths）

```go
viper.OnConfigChange(func(e fsnotify.Event) {
	fmt.Println("Config file changed:", e.Name)
})
viper.WatchConfig()
```

#### 自定义配置源

虽然Viper自带多种配置源，但是这也不妨碍我们自定义。

示例：

```go
viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

// 自定义案例
var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

viper.ReadConfig(bytes.NewBuffer(yamlExample))

viper.Get("name") // 这里将输出 "steve"
```

#### 从Viper中获取配置值

在Viper中有多个获取配置文件内容的方法，可根据值的类型进行选择：

- `Get(key string) : interface{}`
- `GetBool(key string) : bool`
- `GetFloat64(key string) : float64`
- `GetInt(key string) : int`
- `GetIntSlice(key string) : []int`
- `GetString(key string) : string`
- `GetStringMap(key string) : map[string]interface{}`
- `GetStringMapString(key string) : map[string]string`
- `GetStringSlice(key string) : []string`
- `GetTime(key string) : time.Time`
- `GetDuration(key string) : time.Duration`
- `IsSet(key string) : bool`
- `AllSettings() : map[string]interface{}`

注：所有GET方法在没有找到对应的配置参数时将返回0，判断对应的配置参数是否存在可用IsSet()函数

示例：

```go
viper.GetString("logfile") // 不区分大小写
if viper.GetBool("verbose") {
	fmt.Println("verbose enabled")
}
```

#### 访问嵌套类型的配置文件

Viper默认是支持访问嵌套类型的，例如一下JSON，如果要访问其中的某个被嵌套的Key可如此操作：

```go
{
    "host": {
        "address": "localhost",
        "port": 5799
    },
    "datastore": {
        "metric": {
            "host": "127.0.0.1",
            "port": 3099
        },
        "warehouse": {
            "host": "198.0.0.1",
            "port": 2112
        }
    }
}
```

```go
GetString("datastore.metric.host") // (返回 "127.0.0.1")
```

可以通过"."进行嵌套key的访问/获取。

如果你有个嵌套key是一个数组，那我们可以直接读取其下标进行访问：

```go
{
    "host": {
        "address": "localhost",
        "ports": [
            5799,
            6029
        ]
    },
    "datastore": {
        "metric": {
            "host": "127.0.0.1",
            "port": 3099
        },
        "warehouse": {
            "host": "198.0.0.1",
            "port": 2112
        }
    }
}
```

```go
GetInt("host.ports.1") // 返回 6029
```

还有一种极端情况，就是你某个key的名称中就带"."，那么有人就会疑惑会不会无法操作读取呀？放心，Viper已经想到这个问题了，如果有对应的路径则直接取，没有的话才会去嵌套取值：

示例：

```go
{
    "datastore.metric.host": "0.0.0.0",
    "host": {
        "address": "localhost",
        "port": 5799
    },
    "datastore": {
        "metric": {
            "host": "127.0.0.1",
            "port": 3099
        },
        "warehouse": {
            "host": "198.0.0.1",
            "port": 2112
        }
    }
}
```

```go
GetString("datastore.metric.host") // 返回 "0.0.0.0"
```

看，这里返回的是"0.0.0.0" 而不是“127.0.0.1”，我称之为精准匹配原则（我自己定义的哈哈哈哈哈）

#### 以String的形式返回所有配置内容

很多时候我们需要看某个应用程序的配置时是在服务器运行期间，这个时候我们并不是想把他的配置都写入文件中而仅仅是单纯的看一下，那么这个时候AllSettings()这个函数就排上用场了。

```go
import (
	yaml "gopkg.in/yaml.v2"
	// ...
)

func yamlStringSettings() string {
	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}
```

AllSettings()会以String类型返回当前配置中的所有配置内容，

#### 多个Viper实例的使用

上面我们说过：

> Viper虽然支持多配置同时使用，但是一个Viper实例只能寻一个配置路径。

所有如果想支持多配置启用，那么你就多实现几个Viper的实例就可以了。

```go
x := viper.New()
y := viper.New()

x.SetDefault("ContentDir", "content")
y.SetDefault("ContentDir", "foobar")
```

OK，Just this ，See you next~