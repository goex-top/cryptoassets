package main

import (
	"github.com/jinzhu/gorm"
	"github.com/nntaoli-project/GoEx"
)

type Asset struct {
	gorm.Model
	AccountId string
	Btc       float64
	Usdt      float64
}

type CoinAsset struct {
	gorm.Model
	AccountId    string
	CoinName     string
	Amount       float64
	FrozenAmount float64
}

type AccountApi struct {
	goex.APIConfig
	AccountId string
}

type Config struct {
	Freq int          `toml:"freq"`
	Apis []AccountApi `toml:"api"`
}
