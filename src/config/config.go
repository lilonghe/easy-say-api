package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type config struct {
	DB *xorm.Engine

	GithubAppId     string
	GithubAppSecret string

	MySQLName string `json:"mysql_name"`
	MySQLPass string `json:"mysql_pass"`
}

var (
	Global = config{}
)

func Init() {
	initUserConfig()
	initDB()
}

func initDB() {
	engine, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@/easy-say?charset=utf8", Global.MySQLName, Global.MySQLPass))
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

func initUserConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetConfigName("config.json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config.json not found")
		} else {
			panic(err)
		}
	}

	Global.GithubAppId = viper.GetString("github_app_id")
	Global.GithubAppSecret = viper.GetString("github_app_secret")
	Global.MySQLName = viper.GetString("mysql_name")
	Global.MySQLPass = viper.GetString("mysql_pass")
}
