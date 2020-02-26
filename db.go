package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type OrmManager struct {
	*gorm.DB
}

func initOrm() *gorm.DB {
	orm, err := gorm.Open("sqlite3", "./assets.db3")
	if err != nil {
		panic(err)
	}
	orm.AutoMigrate(Account{}, Asset{}, CoinAsset{})
	orm.LogMode(true)
	return orm
}

func (om *OrmManager) AddAccount(account Account) error {
	return om.Create(&account).Error
}

func (om *OrmManager) HasNickName(nickname string) bool {
	acc := om.FindAccountFromNickName(nickname)
	if acc.ID == 0 {
		return false
	}
	return true
}

func (om *OrmManager) FindAccountFromNickName(nickname string) Account {
	acc := Account{NickName: nickname}
	om.Where("nick_name = ?", nickname).First(&acc)
	return acc
}

func (om *OrmManager) GetAssetsFromNickname(nickname string) []Asset {
	acc := om.FindAccountFromNickName(nickname)
	assets := make([]Asset, 0)
	om.Where("nick_name = ?", nickname).First(&acc)
	om.Model(&acc).Related(&assets)
	return assets
}

func (om *OrmManager) GetCoinsFromAssetId(id uint) []CoinAsset {
	asset := Asset{}
	asset.ID = id
	as := make([]CoinAsset, 0)
	om.Model(&asset).Related(&as)
	return as
}
