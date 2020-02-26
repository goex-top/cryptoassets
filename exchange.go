package main

import (
	"github.com/goex-top/market_center"
	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
)

type Exchange struct {
	name     string
	nickname string
	spot     goex.API
	future   []goex.FutureRestAPI
}

func initExchanges(config Config) {
	orm.Find(&accounts)

	for _, v := range accounts {
		secKey := AESECBDecrypt([]byte(v.ApiSecretKey), []byte(config.User.Password))
		passKey := AESECBDecrypt([]byte(v.ApiPassphrase), []byte(config.User.Password))
		exc, isOk := List[v.ExchangeName]
		if !isOk {
			continue
		}
		exchange := Exchange{name: v.ExchangeName, nickname: v.NickName}
		for _, name := range exc {
			if market_center.IsFutureExchange(name) {
				if string(passKey) != "" {
					exchange.future = append(exchange.future,
						builder.NewAPIBuilder().HttpProxy(config.Proxy).
							APIKey(v.ApiKey).APISecretkey(string(secKey)).ApiPassphrase(string(passKey)).
							BuildFuture(name))
				} else {
					exchange.future = append(exchange.future,
						builder.NewAPIBuilder().HttpProxy(config.Proxy).
							APIKey(v.ApiKey).APISecretkey(string(secKey)).
							BuildFuture(name))
				}
			} else {
				if string(passKey) != "" {
					exchange.spot = builder.NewAPIBuilder().HttpProxy(config.Proxy).
						APIKey(v.ApiKey).APISecretkey(string(secKey)).ApiPassphrase(string(passKey)).Build(name)
				} else {
					exchange.spot = builder.NewAPIBuilder().HttpProxy(config.Proxy).
						APIKey(v.ApiKey).APISecretkey(string(secKey)).Build(name)
				}
			}
		}
		exchanges = append(exchanges, exchange)
	}
}
