package main

import (
	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initOrm() *gorm.DB {
	orm, err := gorm.Open("sqlite3", "./assets.db3")
	if err != nil {
		panic(err)
	}
	orm.AutoMigrate(Asset{}, CoinAsset{})
	return orm
}

func loadConfig() (Config, error) {
	var conf Config
	_, err := toml.DecodeFile("account_api.toml", &conf)
	return conf, err
}

func main() {
	initOrm()
	loadConfig()
}
