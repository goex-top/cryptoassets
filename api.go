package main

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func GetSettings(c echo.Context) error {
	type AccountSummary struct {
		gorm.Model
		NickName     string
		ExchangeName string
		ApiKey       string
	}
	acc := make([]AccountSummary, 0)
	for _, v := range accounts {
		acc = append(acc, AccountSummary{
			Model:        v.Model,
			NickName:     v.NickName,
			ExchangeName: v.ExchangeName,
			ApiKey:       v.ApiKey,
		})
	}
	return c.JSON(http.StatusOK, acc)
}

func AddSettings(c echo.Context) error {
	nick_name := c.FormValue("nick_name")
	exchange_name := c.FormValue("exchange_name")
	api_key := c.FormValue("api_key")
	sec_key := c.FormValue("sec_key")
	pass_key := c.FormValue("pass_key")

	if orm.HasNickName(nick_name) {
		return c.String(http.StatusOK, `{code:3000,error: "昵称重名, 请检查"}`)
	}
	account := Account{
		NickName:      nick_name,
		ExchangeName:  exchange_name,
		ApiKey:        api_key,
		ApiSecretKey:  sec_key,
		ApiPassphrase: pass_key,
	}
	orm.AddAccount(account)
	accounts = append(accounts, account)

	return c.String(http.StatusOK, "{}")
}

func GetAssetHistory(c echo.Context) error {
	all := make([][]Asset, 0)
	maxSize := 0
	for _, v := range accounts {
		acc := orm.GetAssetsFromNickname(v.NickName)
		all = append(all, acc)
		if len(acc) > maxSize {
			maxSize = len(acc)
		}
	}
	history := make([]Asset, 0)
	for i := 0; i < maxSize; i++ {
		total := Asset{}
		for _, v := range all {
			index := i - (maxSize - len(v))
			if index < 0 {
				continue
			}
			total.Model = v[index].Model
			total.Btc += v[index].Btc
			total.Btc_Usdt += v[index].Btc_Usdt
			total.Usdt += v[index].Usdt
			total.Usdt_Usd += v[index].Usdt_Usd
			total.Usd_Cny += v[index].Usd_Cny
		}
		history = append(history, Asset{
			Model:    total.Model,
			Btc:      total.Btc,
			Usdt:     total.Usdt,
			Btc_Usdt: total.Usdt,
			Usdt_Usd: total.Usdt_Usd,
			Usd_Cny:  total.Usd_Cny,
		})
	}
	return c.JSON(http.StatusOK, history)
}

func GetCurrentAsset(c echo.Context) error {
	type NicknameAsset struct {
		ID       uint
		NickName string
		Btc      float64
		Usdt     float64
	}
	all := make([]NicknameAsset, 0)
	for _, v := range accounts {
		acc := orm.GetAssetsFromNickname(v.NickName)
		all = append(all, NicknameAsset{
			ID:       acc[len(acc)-1].ID,
			NickName: v.NickName,
			Btc:      acc[len(acc)-1].Btc,
			Usdt:     acc[len(acc)-1].Usdt,
		})
	}

	return c.JSON(http.StatusOK, all)
}

func GetCurrentCoins(c echo.Context) error {
	//nick_name := c.QueryParam("nick_name")
	ids := c.QueryParam("ID")
	id, _ := strconv.Atoi(ids)
	return c.JSON(http.StatusOK, orm.GetCoinsFromAssetId(uint(id)))
}
