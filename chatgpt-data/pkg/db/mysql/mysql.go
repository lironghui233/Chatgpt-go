package mysql

import (
	"chatgpt-data/pkg/config"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitMysql() {
	var err error
	cnf := config.GetConf()
	if cnf.Mysql.DSN == "" {
		panic("数据库连接字符串不能为空")
	}
	db, err = sql.Open("mysql", cnf.Mysql.DSN)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(cnf.Mysql.MaxOpenConn)
	db.SetMaxIdleConns(cnf.Mysql.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(cnf.Mysql.MaxLifeTime))
}

func GetDB() *sql.DB {
	return db
}
