package mysql

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db      *sql.DB
	dbMutex sync.Mutex
)

// ---------------------------------------------------------------------------------------------------------------------

// 初始化 MySQL，仅支持一个实例
func Init(dataSource string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	tmpDB, err := sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}

	tmpDB.SetMaxOpenConns(1000)
	tmpDB.SetMaxIdleConns(200)

	if err := tmpDB.Ping(); err != nil {
		_ = tmpDB.Close()
		return err
	} else {
		db = tmpDB
	}

	return nil
}

// 获取DB
func GetModel(tx ...*sql.Tx) *Model {
	if len(tx) == 0 {
		return &Model{}
	} else {
		return &Model{Tx: tx[0]}
	}
}

// 事务开始
func Begin() (*sql.Tx, error) {
	return db.Begin()
}

// 释放资源
func FreeDB() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db != nil {
		_ = db.Close()
	}
}
