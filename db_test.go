package main

import (
	"testing"
)

var (
	db = initOrm()
)

func init() {
}
func TestOrmCreate(t *testing.T) {
	asset := Asset{
		Btc:  1,
		Usdt: 2222,
		Coins: []CoinAsset{{
			CoinName:     "BTC",
			Amount:       1000,
			FrozenAmount: 20000,
		}, {
			CoinName:     "ETH",
			Amount:       10,
			FrozenAmount: 20,
		}},
	}
	db.Create(&asset)
	db.Save(&asset)
}

func TestOrmRelate(t *testing.T) {
	asset := Asset{}
	//asset.ID = 6
	coins := make([]CoinAsset, 0)
	db.Find(&asset, 5)
	db.Model(&asset).Related(&coins)
	t.Log(asset)
	t.Log(coins)
}
