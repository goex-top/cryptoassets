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
