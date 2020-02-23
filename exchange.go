package main

import (
	"github.com/goex-top/market_center"
	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
)

type Exchange struct {
	spot   goex.API
	future goex.FutureRestAPI
}

func initExchanges(config Config) {
	orm.Find(&accounts)

	for _, v := range accounts {
		secKey := AESECBDecrypt([]byte(v.ApiSecretKey), []byte(config.User.Password))
		//passKey := AESECBDecrypt([]byte(v.ApiPassphrase), []byte(config.User.Password))
		if market_center.IsFutureExchange(v.ExchangeName) {
			exchanges = append(exchanges, Exchange{
				spot:   nil,
				future: builder.NewAPIBuilder().HttpProxy(config.Proxy).APIKey(v.ApiKey).APISecretkey(string(secKey)).BuildFuture(v.ExchangeName),
			})
		} else {
			exchanges = append(exchanges, Exchange{
				spot:   builder.NewAPIBuilder().HttpProxy(config.Proxy).APIKey(v.ApiKey).APISecretkey(string(secKey)).Build(v.ExchangeName),
				future: nil,
			})
		}
	}
}
