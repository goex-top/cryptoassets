package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initOrm() *gorm.DB {
	orm, err := gorm.Open("sqlite3", "./assets.db3")
	if err != nil {
		panic(err)
	}
	orm.AutoMigrate(Account{}, Asset{}, CoinAsset{})
	orm.LogMode(true)
	return orm
}

func addAccount(account Account) {
	orm.Create(&account)
}

func getAssetsFromNickname(nickname string) []Asset {
	acc := Account{NickName: nickname}
	assets := make([]Asset, 0)
	orm.First(&acc)
	orm.Model(&acc).Related(&assets)
	return assets
}

func getCoinsFromAssetId(id uint) []CoinAsset {
	asset := Asset{}
	asset.ID = id
	as := make([]CoinAsset, 0)
	orm.Model(&asset).Related(&as)
	return as
}
