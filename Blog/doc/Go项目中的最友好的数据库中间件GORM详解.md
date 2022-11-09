## <a id="6">GORM-Go语言中号称最友好的ORM</a>

ORM即Object-Relationl Mapping，就是对象关系映射。它的作用就是在对象和关系数据库中间做一个映射，这样我们在操作数据库的时候就不需要直接去和SQL打架。

其在Go圈子中非常火，功能覆盖范围也非常广，这节我们不做太深入的研究，先搞明白怎么用即可。

#### 安装

```go
go get -u gorm.io/gorm
```

如果上面这个你用不了就用下面这个

```go
go get github.com/jinzhu/gorm
```

#### 快速开始案例---Sqlite数据库链接

```go
package main
//导包
import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)
//定义一个模型
type Product struct {
  gorm.Model
  Code  string
  Price uint
}

func main() {
  //根据配置文件链接名为test.db的数据库
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // 在数据库中创建对应的Product表
  db.AutoMigrate(&Product{})

  // 创建一条数据
  db.Create(&Product{Code: "D42", Price: 100})

  // 查询数据
  var product Product
  db.First(&product, 1) 
  db.First(&product, "code = ?", "D42") 

  // 更新数据
  db.Model(&product).Update("Price", 200)
  // 更新多条数据
  db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
  db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

  // 删除数据
  db.Delete(&product, 1)
}
```

#### 快速开始案例---MySQL数据库链接

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

GORM还为MySQL Driver提供了一些可以在初始化期间使用的高级配置，如下：

```go
db, err := gorm.Open(mysql.New(mysql.Config{
  DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // 数据源
  DefaultStringSize: 256, // 字符串的默认大小
  DisableDatetimePrecision: true, // 在MySQL 5.6版本之前不支持日期时间精度，这里给禁用
  DontSupportRenameIndex: true, // 在MySQL 5.7之前不支持重命名索引
  DontSupportRenameColumn: true, // MySQL 8不支持重命名列
  SkipInitializeWithVersion: false, // 根据当前MySQL版本自动配置
}), &gorm.Config{})
```

#### 自定义Driver

在GORM中允许使用DriverName选项去自定义MySQL Driver，如下：

```go
import (
  _ "example.com/my_mysql_driver"
  "gorm.io/gorm"
)

db, err := gorm.Open(mysql.New(mysql.Config{
  DriverName: "my_mysql_driver",
  DSN: "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // 数据源
}), &gorm.Config{})
```

GROM允许去初始化*gorm.DB 去持有一个现有的数据库链接，如下：

```go
import (
  "database/sql"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

sqlDB, err := sql.Open("mysql", "mydb_dsn")
gormDB, err := gorm.Open(mysql.New(mysql.Config{
  Conn: sqlDB,
}), &gorm.Config{})
```

除了以上举的这两个数据库外，其实还有PostgreSQL、SQLServer、Clickhouse、Pool 等，这些就自行去查阅吧。

#### 创建模型

模型就是采用Go的基本类型、指针/别名或者其他用户自定义类型等信息去实现Scanner和Valuer接口。

直接看示例：

```go
type User struct {
  ID           uint
  Name         string
  Email        *string
  Age          uint8
  Birthday     *time.Time
  MemberNumber sql.NullString
  ActivatedAt  sql.NullTime
  CreatedAt    time.Time
  UpdatedAt    time.Time
}
```

在GORM中有一个默认的约定，以模型创建的表默认表名就是结构体名称，默认字段就是结构体中的字段，默认的primary key就是ID，并且自带CreateAt和UpdateAt字段去记录创建和更新时间。如果GROM的这个默认约定不满足你的使用，那么你也可以抛弃这个约定去自行定义。

#### gorm.Model

这个在我们开源十年项目中就用过，且也给大家讲述过，它是GROM默认的结构体，你可以直接把它嵌入进自己的结构体中，然后你的结构体将包含它所自带的几个字段：ID，CreateAt，UpdateAt,DeletedAt

```go
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

#### 字段的权限等级

使用GROM去做增删查改的时候对所有字段是拥有所有权限的，其实GROM是允许我们通过Tag去修改字段的等级权限的，所以在正常的开发中我们可以去自定义字段的只读、只写、只允许创建、只允许更新或者忽略等权限。

注：这里需要注意的一点是，如果某个字段设置了Tag为忽略，那么在我们去创建这个表的时候也是会忽略这个字段的。

具体设置参考下表：

```go
type User struct {
  Name string `gorm:"<-:create"` // 允许读和创建
  Name string `gorm:"<-:update"` // 允许读和更新
  Name string `gorm:"<-"`        // 允许读写（创建和更新）
  Name string `gorm:"<-:false"`  // 允许读 ，不允许写
  Name string `gorm:"->"`        // 只读（除非有别的特殊配置不然就是禁用写）
  Name string `gorm:"->;<-:create"` // 允许读和创建
  Name string `gorm:"->:false;<-:create"` // 只是创建（禁用从数据库的读取权限）
  Name string `gorm:"-"`            // 在结构体读写时忽略此字段
  Name string `gorm:"-:all"`        // 在结构体读写、甚至是迁移表时忽略此字段
  Name string `gorm:"-:migration"`  // 在结构体迁移表时忽略该字段
}
```

#### 默认CreatAt/UpdateAt和默认时间格式的使用

在GROM中CreatAt和UpdateAt是默认的，且会在我们操作相关数据是自动记录并更新其时间内容。如果我们想自定义同类型的字段只需要加上autoCreateTime和autoUpdateTime的Tag即可。同时GROM默认的时间格式是time.Time,如果你想用秒或者其它类型去记录，只需要把time/Time更改为Int或者对应的其它类型即可。

示例：

```go
type User struct {
  CreatedAt time.Time // 创建时不指定具体指会默认设置为当前时间
  UpdatedAt int       // 在更新时设置为当前unix秒类型的值，创建时此值为0
  Updated   int64 `gorm:"autoUpdateTime:nano"` // 使用unix nano类型为更新时间
  Updated   int64 `gorm:"autoUpdateTime:milli"`// 使用 unix milli 类型为更新时间
  Created   int64 `gorm:"autoCreateTime"`      // 使用unix为创建时间
}
```

#### 嵌入式结构

我喜欢称之为集成式结构，客官们随意哈。对于此结构我们可以实现更多种类型和更复杂的结构，并且便于维护和读写。示例：

```go
type Author struct {
  Name  string
  Email string
}

type Blog struct {
  ID      int
  Author  Author `gorm:"embedded"`
  Upvotes int32
}
// 上面那种写法与下面的这种写法是一样的效果
type Blog struct {
  ID    int64
  Name  string
  Email string
  Upvotes  int32
}
```

同时可以使用embeddedPrefix这个标签为字段在数据库中添加前缀：

```go
type Blog struct {
  ID      int
  Author  Author `gorm:"embedded;embeddedPrefix:author_"`
  Upvotes int32
}
// equals
type Blog struct {
  ID          int64
  AuthorName  string
  AuthorEmail string
  Upvotes     int32
}
```

这么看的话就很方便了是吧。

在GROM中具体有哪些Tag，我这里就不去一一列举了，就像字典一样，你用到了再去查效果会好很多。

```go
column
//指定 db 列名

type
//列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 bool、int、uint、float、string、time、bytes 并且可以和其他标签一起使用，例如：not null、size, autoIncrement… 像 varbinary(8) 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：MEDIUMINT UNSIGNED not NULL AUTO_INSTREMENT


size
//指定列大小，例如：size:256


primaryKey
//指定列为主键


unique
//指定列为唯一


default
//指定列的默认值


precision
//指定列的精度


scale
//指定列大小


not null
//指定列为 NOT NULL


autoIncrement
//指定列为自动增长


embedded
//嵌套字段


embeddedPrefix
//嵌入字段的列名前缀


autoCreateTime
//创建时追踪当前时间，对于 int 字段，它会追踪时间戳秒数，您可以使用 nano/milli 来追踪纳秒、毫秒时间戳，例如：autoCreateTime:nano


autoUpdateTime
//创建 / 更新时追踪当前时间，对于 int 字段，它会追踪时间戳秒数，您可以使用 nano/milli 来追踪纳秒、毫秒时间戳，例如：autoUpdateTime:milli


index
//根据参数创建索引，多个字段使用相同的名称则创建复合索引，查看 索引 获取详情


uniqueIndex
//与 index 相同，但创建的是唯一索引


check
//创建检查约束，例如 check:age > 13，查看 约束 获取详情


<-
//设置字段写入的权限， <-:create 只创建、<-:update 只更新、<-:false 无写入权限、<- 创建和更新权限


->
//设置字段读的权限，->:false 无读权限

-
//忽略该字段，- 无读写权限
```

### 增删查改。

在学习GORM的时候，除了上面这些东西是需要我们理解之外，以下内容才是我们需要学的重点。

##### Create

```go
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

result := db.Create(&user) // 将数据指针传给Create

user.ID             // 返回数据的ID prinary key 
result.Error        // 返回错误信息
result.RowsAffected // 返回插入的记录数
//创建一条记录并为指定字段赋值
db.Select("Name", "Age", "CreatedAt").Create(&user)//INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775")
//创建一条记录并忽略掉给省略字段传值
db.Omit("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775")

//批量插入
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)

for _, user := range users {
  user.ID // 1,2,3
}
//也可以使用CreateInBatches指定批量的范围
var users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}

// batch size 100
db.CreateInBatches(users, 100)

```

##### 检索单个对象

```go
// 获取排在第一的值（以prinary key排序）
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// 随便获取一个值
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// 获取排在最后一个的值（以prinary key排序）
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

result := db.First(&user)
result.RowsAffected // 返回一共查到了多少条记录
result.Error        // 返回错误信息

// 检查是否存在错误 ErrRecordNotFound
errors.Is(result.Error, gorm.ErrRecordNotFound)

```

First和Last这两个只有在目标结构的指针做参数时才起作用，此外，如果这个模型没有primary key那返回的顺序就是按照第一个字段排序的结果去取值。

按照指定主键去检索

```go
db.First(&user, 10)
// SELECT * FROM users WHERE id = 10;

db.First(&user, "10")
// SELECT * FROM users WHERE id = 10;

db.Find(&users, []int{1,2,3})
// SELECT * FROM users WHERE id IN (1,2,3);
//如果主键是字符串（例如，像 uuid），则查询如下：
db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";

```

注：在主键是字符串时，请谨防SQL注入安全事故。

```go
//检索所有对象
// 获取所有记录
result := db.Find(&users)
// SELECT * FROM users;

```

条件查询

```go
// 获取第一条匹配的记录
db.Where("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

// 获取所有匹配的记录
db.Where("name <> ?", "jinzhu").Find(&users)
// SELECT * FROM users WHERE name <> 'jinzhu';


//以下的IN LIKE AND Time等就跟我们学的SQL语句中的用法是一样的了
// IN
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

// LIKE
db.Where("name LIKE ?", "%jin%").Find(&users)
// SELECT * FROM users WHERE name LIKE '%jin%';

// AND
db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

// Time
db.Where("updated_at > ?", lastWeek).Find(&users)
// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';

// BETWEEN
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';

```

暂时查询就说这么多吧，够现阶段用了。后面还有排序、限制条件、分组、Join等等跟我们之前语言中遇到的基本是一致的。

##### Update

```go
//Save将会把所有字段都存储
db.First(&user)

user.Name = "jinzhu 2"
user.Age = 100
db.Save(&user)
// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;


//更新单列
// 按条件更新
db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

// 更新指定条件的记录
db.Model(&user).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

//使用条件+模型值更新
db.Model(&user).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;


//更新多列
// 使用结构模型更新属性，但是只会更新非0的字段
db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// 使用map更新属性
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;


//更新指定的属性

db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;

db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// 注意肯定是非0值
db.Model(&user).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})
// UPDATE users SET name='new_name', age=0 WHERE id=111;

// 选择所有字段包括非0
db.Model(&user).Select("*").Update(User{Name: "jinzhu", Role: "admin", Age: 0})

// 选择所有字段包括非0 但是忽略“Role”
db.Model(&user).Select("*").Omit("Role").Update(User{Name: "jinzhu", Role: "admin", Age: 0})
```

##### Delete

```go
//删除单条记录
db.Delete(&email)
// DELETE from emails where id = 10;

// 按条件删除
db.Where("name = ?", "jinzhu").Delete(&email)
// DELETE from emails where id = 10 AND name = "jinzhu";

//按primary key去删除
db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;

db.Delete(&User{}, "10")
// DELETE FROM users WHERE id = 10;

db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);

//批量删除
db.Where("email LIKE ?", "%jinzhu%").Delete(&Email{})
// DELETE from emails where email LIKE "%jinzhu%";

db.Delete(&Email{}, "email LIKE ?", "%jinzhu%")
// DELETE from emails where email LIKE "%jinzhu%";

```

好了 ，就先分享这么多吧。其实Gorm的内容非常多，但是其中的很多内容都是百变不离其宗。用到了再来查就是了，目前这些对于我们现阶段的项目来说够用就行。

OK .Just this,See you next...