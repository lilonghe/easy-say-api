package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type config struct {
	DB *xorm.Engine
}

var (
	Global = config{}
)

func Init() {
	initDB()
}

func initDB() {
	engine, err := xorm.NewEngine("mysql", "root:lilonghe@/easy-say?charset=utf8")
	if err != nil {
		panic(err)
	}
	err = engine.Ping()
	if err != nil {
		panic(err)
	}
	engine.ShowSQL(true)
	engine.Logger().SetLevel(log.LOG_DEBUG)
	Global.DB = engine
	fmt.Println("init db ok!")
}
