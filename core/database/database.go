package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
	"zxi_go/core/config"
)

func MySqlInit() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.MySqlConf["user"],
		config.MySqlConf["pwd"],
		config.MySqlConf["type"],
		config.MySqlConf["address"],
		config.MySqlConf["port"],
		config.MySqlConf["database"],
	)
	LocalDB, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Failed to connect to mysql,err:" + err.Error())
		os.Exit(1)
	}
	LocalDB.DB().SetMaxIdleConns(10)
	LocalDB.DB().SetMaxOpenConns(100)
	return LocalDB
}
