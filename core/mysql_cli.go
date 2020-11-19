package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
)

var LocalDB *gorm.DB

func MySqlInit() {
	var err error
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		MySqlConf["user"],
		MySqlConf["pwd"],
		MySqlConf["type"],
		MySqlConf["address"],
		MySqlConf["port"],
		MySqlConf["database"],
	)
	LocalDB, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Failed to connect to mysql,err:" + err.Error())
		os.Exit(1)
	}
	LocalDB.DB().SetMaxIdleConns(10)
	LocalDB.DB().SetMaxOpenConns(100)
}
