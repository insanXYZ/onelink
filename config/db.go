package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewDb() *sql.DB {
	sql, err := sql.Open("mysql", "root:@tcp(localhost:3306)/rad?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return sql
}
