package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"zxi_network_disk_go/conf"
)

var Conn *sql.DB

func MySqlInit() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",
		conf.MySqlConf["user"],
		conf.MySqlConf["pwd"],
		conf.MySqlConf["type"],
		conf.MySqlConf["address"],
		conf.MySqlConf["port"],
		conf.MySqlConf["database"],
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
