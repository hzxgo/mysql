# mysql
以前使用了别人的MySQL ORM，工作中用起来确实方便而且省事，但是也会觉得有些ORM太重，且有些地方不太好用，而且别人的ORM用久了之后会发现自己的SQL水平基本停滞不前，反而会淡化，故自己写了个MySQL ORM，自己平时基本就是使用该ORM。

### 使用方法
```
import (
    "os"

    "github.com/hzxgo/log"
    "github.com/hzxgo/mysql"
)

func init() {

    // init mysql
    dataSource := "write your mysql data source"
    if err := mysql.Init(dataSource); err != nil {
        log.Errorf("init mysql failed | %v", err)
        os.Exit(-1)
    }
}

type User struct {
	mysql.Model `db:"-"`  // db 中没有该字段，故不会对该字段进行插入和更新
	ID          int64  `db:"ID auto_increment"` // db 对应字段名称为“ID”，且为自动递增字段
	Username    string  `db:"Username"`          // db 对应字段名称为“Username”
	Password    string           // db 对应字段名称默认就是结构体的字段名称
	CreateTime  int64            // db 对应字段名称默认就是结构体的字段名称
	Describe    mysql.NullString // 描述：默认为NULL类型
	Address     string  `db:"-"` // db 中没有该字段，故不会对该字段进行插入和更新
}

func NewUser() *User {
	return &User{
		Model: mysql.Model{
			TableName: "user",
		},
	}
}
```

### insert使用方法
```
# 指定列插入数据
user := NewUser()
var describe mysql.NullString
describe.String = "this is comment"
params := map[string]interface{}{
    "Username":   "hezhixiong",
    "Password":   "1234567890",
    "CreateTime": time.Now().Unix(),
    // "Describe":   describe, // 你也可以注释改行，即只插入以上三个字段的值
}
id, err := user.Insert(params)
if err != nil {
    log.Errorf("insert user failed | %v", err)
    return
} else {
    log.Infof("insert user success, id: %d", id)
}

# 基于对象进行插入数据
u := NewUser()
u.Username = "new_username"
u.Password = "new_password"
id, err := u.Insert(u)
if err != nil {
    log.Errorf("insert user failed | %v", err)
    return
} else {
    log.Infof("insert user success, id: %d", id)
}

# 当然也支持批量插入操作（文档后期再来补）

```
