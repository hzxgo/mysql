package mysql

import (
	"database/sql"
)

// ---------------------------------------------------------------------------------------------------------------------

const (
	maxBatchLimit = 500 // 最大批量操作量
)

const (
	dbTagEmpty         = ""               // 空
	dbTagDiscard       = "-"              // 丢弃
	dbTagAutoIncrement = "auto_increment" // 自动递增
)

// ---------------------------------------------------------------------------------------------------------------------

// NullString is a type that can be null or a string
type NullString struct {
	sql.NullString
}

// NullFloat64 is a type that can be null or a float64
type NullFloat64 struct {
	sql.NullFloat64
}

// NullInt64 is a type that can be null or an int
type NullInt64 struct {
	sql.NullInt64
}

// NullBool is a type that can be null or a bool
type NullBool struct {
	sql.NullBool
}
