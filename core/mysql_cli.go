package core

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var Conn *sql.DB

func MySqlInit() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		MySqlConf["user"],
		MySqlConf["pwd"],
		MySqlConf["type"],
		MySqlConf["address"],
		MySqlConf["port"],
		MySqlConf["database"],
	)
	Conn, _ = sql.Open("mysql", dsn)
	Conn.SetMaxOpenConns(50)
	Conn.SetMaxIdleConns(10)
	err := Conn.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql,err:" + err.Error())
		os.Exit(1)
	}
}
