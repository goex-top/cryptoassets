package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type OrmManager struct {
	*gorm.DB
}

func initOrm(debug bool) *gorm.DB {
	orm, err := gorm.Open("sqlite3", "./assets.db3")
	if err != nil {
		panic(err)
	}
	orm.AutoMigrate(Account{}, Asset{}, CoinAsset{})
	if debug {
		orm.LogMode(true)
	}
	return orm
}

func (om *OrmManager) AddAccount(account Account) (Account, error) {
	o := om.Create(&account)
	return account, o.Error
}

func (om *OrmManager) DeleteAccount(id uint) error {
	acc := Account{}
	acc.ID = id
	o := om.Delete(&acc)
	return o.Error
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

func (om *OrmManager) GetAssetsFromAccountId(id uint) []Asset {
	acc := Account{}
	acc.ID = id
	assets := make([]Asset, 0)
	om.Model(&acc).Related(&assets)
	return assets
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

func (om *OrmManager) AddAsset(asset Asset) Asset {
	om.Create(&asset)
	return asset
}

func (om *OrmManager) AddCoinAssets(coinAssets []CoinAsset) error {
	for _, ca := range coinAssets {
		db := om.Create(&ca)
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (om *OrmManager) AddCoinAsset(coinAsset CoinAsset) CoinAsset {
	om.Create(&coinAsset)
	return coinAsset
}
