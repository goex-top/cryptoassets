package main

import (
	"testing"
)

var (
	db OrmManager
)

func init() {
	db.DB = initOrm()
}

func TestOrmManager_AddAccount(t *testing.T) {
	t.Log(db.AddAccount(Account{
		NickName:      "aaa",
		ExchangeName:  "",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
	}))
	t.Log(db.AddAccount(Account{
		NickName:      "bbb",
		ExchangeName:  "",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
	}))
}

func TestFindAccountFromNickName(t *testing.T) {
	t.Log(db.FindAccountFromNickName("aaa"))
	t.Log(db.FindAccountFromNickName("bbb"))
	t.Log(db.FindAccountFromNickName("ccc"))
}

func TestOrmManager_AddAsset(t *testing.T) {
	t.Log(db.AddAsset(Asset{
		AccountID: 1,
		Btc:       11,
		Usdt:      12,
		Usd:       13,
		Cny:       14,
		Btc_Usdt:  21,
		Btc_Usd:   22,
		Btc_Cny:   23,
		Usdt_Usd:  24,
		Usdt_Cny:  25,
		Usd_Cny:   26,
	}))
}

func TestOrmManager_FindAssets(t *testing.T) {
	t.Log(db.GetAssetsFromAccountId(1))
}
