package main

import (
	"github.com/jinzhu/gorm"
)

// database table
type Account struct {
	gorm.Model
	NickName       string  `gorm:"unique" json:"nick_name"`
	ExchangeName   string  `json:"exchange_name"`
	ApiKey         string  `json:"api_key"`
	ApiSecretKey   string  `json:"api_secret_key"`
	ApiPassphrase  string  `json:"api_passphrase"`
	Assets         []Asset `json:"assets"`
	LastUpdateTime int64   `json:"last_update_time"`
}

type Asset struct {
	gorm.Model
	AccountID uint        `json:"account_id"`
	Btc       float64     `json:"btc"`
	Usdt      float64     `json:"usdt"`
	Usd       float64     `json:"usd"`
	Cny       float64     `json:"cny"`
	Btc_Usdt  float64     `json:"btc_usdt"`
	Btc_Usd   float64     `json:"btc_usd"`
	Btc_Cny   float64     `json:"btc_cny"`
	Usdt_Usd  float64     `json:"usdt_usd"`
	Usdt_Cny  float64     `json:"usdt_cny"`
	Usd_Cny   float64     `json:"usd_cny"`
	Coins     []CoinAsset `json:"coins"`
}

type CoinAsset struct {
	gorm.Model
	AssetID      uint    `json:"asset_id"`
	CoinName     string  `json:"coin_name"`
	Amount       float64 `json:"amount"`
	FrozenAmount float64 `json:"frozen_amount"`
	Btc          float64 `json:"btc"`
	Usdt         float64 `json:"usdt"`
	Usd          float64 `json:"usd"`
	Cny          float64 `json:"cny"`
}

// configure
type Config struct {
	Freq  int    `toml:"freq"`
	Proxy string `toml:proxy`
	User  User   `toml:"user"`
}

type User struct {
	UserName string `json:"username" xml:"username" form:"username" query:"username"`
	Password string `json:"password" xml:"password" form:"password" query:"password"`
}
