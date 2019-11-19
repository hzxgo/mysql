# MySQL ORM
自己写了一个MySQL的ORM，该ORM非常轻量级且功能丰富，自己和团队就是使用该ORM对MySQL进行操作，功能方面支持如下：
* 查询单条记录
* 查询多条记录
* 添加单条记录
* 添加多条记录
* 更新记录
* 删除记录
* 支持事务操作(Tags v1.0.1版本不支持事务)

## 使用介绍
### 初始化MySQL操作
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
		log.Errorf("init mysql | %v", err)
		os.Exit(-1)
	}
}
```

### 模型代码示例
```
import (
	"github.com/hzxgo/mysql"
)

type User struct {
	ID              int64  `db:"ID"`
	Username        string // 若不指定 db 中的名称，则默认为字段名称
	Password        string
	RealName        string
	IsAdmin         int
	State           int
	CreateTime      int64
	UpdateTime      int64
	LatestLoginTime int64
	LoginTimes      int
	Comment         mysql.NullString // db字段类型为：text
	Token           string
}

// ---------------------------------------------------------------------------------------------------------------------

func NewUser() *User {
	return &User{}
}

// 表名
func (u *User) TableName() string {
	return "ddy_user"
}

// 获取单条数据
func (u *User) GetSingleByExp(exp map[string]interface{}) error {

	// 构造查询语句
	builder := mysql.Select("*").Form(u.TableName()).Limit(1)
	rows, err := mysql.SelectWhere(builder, exp)
	if err != nil {
		return err
	}

	// 加载数据
	if err := mysql.LoadStruct(rows, u); err != nil {
		return err
	}

	return nil
}

// 获取单条数据
func (u *User) GetSingleByID(id int64) error {

	// 使用原生SQL查询
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE `ID`=%d LIMIT 1", u.TableName(), id)
	rows, err := mysql.SelectBySql(cmd)
	if err != nil {
		return err
	}

	// 加载数据
	if err := mysql.LoadStruct(rows, u); err != nil {
		return err
	}

	return nil
}

// 获取单条数据
func (u *User) GetLimitByExp(exp map[string]interface{}, offset, limit uint64) ([]*User, error) {

	// 构造查询语句(仅获取指定列数据)
	builder := mysql.Select("`ID`, `Username`, `Password`").Form(u.TableName()).LimitPage(offset, limit)
	rows, err := mysql.SelectWhere(builder, exp)
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, limit)

	// 加载数据
	if count, err := mysql.LoadStructs(rows, &users); err != nil {
		return nil, err
	} else {
		fmt.Println("count: ", count)
	}

	return users, nil
}

// 插入单条数据
// data数据类型：对象指针类型 或 Map 类型
func (u *User) Insert(data ...interface{}) (int64, error) {
	var insertData interface{}

	if len(data) > 0 {
		insertData = data[0]
	} else {
		insertData = u
	}

	return mysql.Insert(u.TableName(), insertData)
}

// 插入多条数据
// data数据类型：对象指针类型 或 Map 类型
func (u *User) MInsert(data ...interface{}) (int64, int64, error) {

	return mysql.MInsert(u.TableName(), data...)
}

// 插入多条数据
func (u *User) BatchInsert(columns []string, params []interface{}) (int64, int64, error) {

	return mysql.BatchInsert(u.TableName(), columns, params)
}

// 更新：基于exp表达式更新data数据
func (u *User) Update(data interface{}, exp interface{}) (int64, error) {

	return mysql.Update(u.TableName(), data, exp)
}

// 删除：基于exp表达式删除数据
func (u *User) Delete(exp interface{}) (int64, error) {

	return mysql.Delete(u.TableName(), exp)
}
```

### 演示代码

```

func main() {

	testSelect()

	testInsert()

	testBatchInsert()

	testUpdate()

	testDelete()
}

func testSelect() {
	user := NewUser()

	// 基于表达式获取单条数据
	exp := map[string]interface{}{
		"ID = ?": 1,
	}
	if err := user.GetSingleByExp(exp); err != nil {
		log.Errorf("User GetSingleByExp err: %v", err)
		return
	}
	log.Infof("User GetSingleByExp: %+v", user)

	// 基于ID获取单条数据
	if err := user.GetSingleByID(2); err != nil {
		log.Errorf("User GetSingleByID err: %v", err)
		return
	}
	log.Infof("User GetSingleByID: %+v", user)

	// 获取翻页数据
	exp = map[string]interface{}{
		"ID > ?": 0,
	}
	users, err := user.GetLimitByExp(exp, 0, 20)
	if err != nil {
		log.Errorf("User GetLimitByExp err: %v", err)
		return
	}
	log.Infof("len(users): %v, users: %+v", len(users), users)
}

func testInsert() {
	user := NewUser()

	// 插入指定列数据，备注：params插入后不会加载至user
	// 以下实际SQL：
	// INSERT INTO `ddy_user` (`Username`,`Password`,`IsAdmin`,`CreateTime`,`UpdateTime`,`LatestLoginTime`,`Comment`)
	//    VALUES('sam','123456','1','1574135084','1574135084','1574135084','')
	timestamp := time.Now().Unix()
	params := map[string]interface{}{
		"Username":        "sam",
		"Password":        "123456",
		"IsAdmin":         1,
		"CreateTime":      timestamp,
		"UpdateTime":      timestamp,
		"LatestLoginTime": timestamp,
		"Comment":         "",
	}
	if id, err := user.Insert(params); err != nil {
		log.Errorf("User Insert err: %v", err)
		return
	} else {
		log.Infof("User Insert Id: %v", id)
		log.Infof("User: %+v", user)
	}

	// 插入对象
	// 以下实际SQL：
	// INSERT INTO `ddy_user` (`IsAdmin`,`UpdateTime`,`ID`,`RealName`,`State`,`CreateTime`,
	//    `LatestLoginTime`,`LoginTimes`,`Comment`,`Token`,`Username`,`Password`)
	//     VALUES('0','0','0','','0','0','0','0','this is test','','test','123')
	user.Username = "test"
	user.Password = "123"
	user.Comment.String = "this is test"
	if id, err := user.Insert(); err != nil {
		log.Errorf("User Insert err: %v", err)
		return
	} else {
		user.ID = id
	}
}

func testBatchInsert() {
	user := NewUser()

	timestamp := time.Now().Unix()
	params1 := map[string]interface{}{
		"Username":        "test1",
		"Password":        "123456",
		"IsAdmin":         1,
		"CreateTime":      timestamp,
		"UpdateTime":      timestamp,
		"LatestLoginTime": timestamp,
		"Comment":         "",
	}
	params2 := map[string]interface{}{
		"Username":        "test2",
		"Password":        "123456",
		"IsAdmin":         1,
		"CreateTime":      timestamp,
		"UpdateTime":      timestamp,
		"LatestLoginTime": timestamp,
		"Comment":         "",
	}

	// 插入两个params
	// 以下实际SQL：
	// INSERT INTO `ddy_user` (`Username`,`Password`,`IsAdmin`,`CreateTime`,`UpdateTime`,`LatestLoginTime`,`Comment`)
	//     VALUES ('test1','123456',1,1574136054,1574136054,1574136054,''),
	//     ('test2','123456',1,1574136054,1574136054,1574136054,'')
	if id, affected, err := user.MInsert(params1, params2); err != nil {
		log.Errorf("User MInsert | %v", err)
		return
	} else {
		log.Infof("User MInsert EffectRows: %v, %v", id, affected)
	}

	user1 := NewUser()
	user1.Username = "test3"
	user1.Password = "123"
	user1.Comment.String = ""

	user2 := NewUser()
	user2.Username = "test4"
	user2.Password = "1234"
	user2.Comment.String = ""

	// 插入两个指针对象
	// 以下实际SQL：
	// INSERT INTO `ddy_user` (`ID`,`Username`,`Password`,`RealName`,`IsAdmin`,`State`,
	//     `CreateTime`,`UpdateTime`,`LatestLoginTime`,`LoginTimes`,`Comment`,`Token`)
	//      VALUES (0,'test3','123','',0,0,0,0,0,0,'',''),(0,'test4','1234','',0,0,0,0,0,0,'','')
	if id, affected, err := user1.MInsert(user1, user2); err != nil {
		log.Errorf("User MInsert | %v", err)
		return
	} else {
		log.Infof("User MInsert EffectRows: %v %v", id, affected)
	}

	// batch insert
	// 以下实际SQL：
	// INSERT INTO `ddy_user` (`Username`,`Password`,`Comment`)
	//   VALUES ('name_0','123',''),('name_1','123',''),('name_2','123','')
	columns := []string{"Username", "Password", "Comment"}
	values := make([]interface{}, 0, 3)
	for i := 0; i < 3; i++ {
		values = append(values, []interface{}{fmt.Sprintf("name_%v", i), "123", ""})
	}
	id, affected, err := user.BatchInsert(columns, values)
	if err != nil {
		log.Errorf("User BatchInsert err: %v", err)
		return
	}
	log.Infof("User BatchInsert id: %v, affected: %v", id, affected)
}

func testUpdate() {
	user := NewUser()

	// 以下实际SQL：
	// UPDATE `ddy_user` SET `Password`='md5(username)'
	params := map[string]interface{}{
		"Password": "md5(username)",
	}
	affected, err := user.Update(params, nil)
	if err != nil {
		log.Errorf("User Update err: %v", err)
		return
	}
	log.Infof("User Update Affected: %v", affected)
}

func testDelete() {
	user := NewUser()

	// 以下实际SQL：
	// DELETE FROM `ddy_user`  WHERE (ID > '5') AND (IsAdmin = '1' OR LoginTimes = '0')
	exp := map[string]map[string]interface{}{
		"AND": {
			"ID > ?": 5,
		},
		"OR": {
			"IsAdmin = ?":    1,
			"LoginTimes = ?": 0,
		},
	}
	affected, err := user.Delete(exp)
	if err != nil {
		log.Errorf("User Delete err: %v", err)
		return
	}
	log.Infof("User Delete Affected: %v", affected)
}
```
