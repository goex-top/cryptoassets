package main

import (
	"context"
	"github.com/goex-top/market_center"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
	"time"
)

type Exchange struct {
	name      string
	nickname  string
	accountId uint
	spot      goex.API
	future    []goex.FutureRestAPI
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
		exchange := Exchange{name: v.ExchangeName, nickname: v.NickName, accountId: v.ID}
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

func UpdateAccounts() {
	usdcny := GetUsdCny()
	btcusd := GetBtcUsd()
	usdtusd := GetUsdtUsd()

	btcusdt := btcusd * usdtusd
	btccny := btcusd * usdcny
	usdtcny := usdtusd * usdcny
	for _, ex := range exchanges {
		coins := make(map[string]CoinAsset)
		if ex.spot != nil {
			acc, err := ex.spot.GetAccount()
			if err != nil {
				continue
			}
			for _, sub := range acc.SubAccounts {
				total := sub.Amount + sub.ForzenAmount
				coin := CoinAsset{
					CoinName:     sub.Currency.String(),
					Amount:       sub.Amount,
					FrozenAmount: sub.ForzenAmount,
				}
				if sub.Currency == goex.USDT {
					coin.Btc = total / btcusdt
					coin.Usdt = total
					coin.Usd = total / usdtusd
					coin.Cny = total / usdtcny
				} else if sub.Currency == goex.BTC {
					coin.Btc = total
					coin.Usdt = total * btcusdt
					coin.Usd = total * btcusd
					coin.Cny = total * btccny
				} else {
					usdt_pair := goex.NewCurrencyPair(sub.Currency, goex.USDT)
					btc_pair := goex.NewCurrencyPair(sub.Currency, goex.BTC)
					usdt_ticker, err := ex.spot.GetTicker(usdt_pair)

					if err != nil {
						btc_ticker, err := ex.spot.GetTicker(btc_pair)
						if err == nil {
							btc := total * btc_ticker.Last
							coin.Btc = btc
							coin.Usdt = btc * btcusdt
							coin.Usd = btc * btcusd
							coin.Cny = btc * btccny
						} else {
							logger.Printf("coin:%s not ticker for value caculate", sub.Currency.String())
						}
					} else {
						usdt := total * usdt_ticker.Last
						coin.Btc = usdt / btcusdt
						coin.Usdt = usdt
						coin.Usd = usdt / usdtusd
						coin.Cny = usdt / usdtcny
					}
				}
				coins[sub.Currency.String()] = coin
			}
		}

		if len(ex.future) > 0 {
			for _, future := range ex.future {
				if future == nil {
					continue
				}
				acc, err := future.GetFutureUserinfo()
				if err != nil {
					continue
				}
				for _, sub := range acc.FutureSubAccounts {
					total := sub.AccountRights // todo: pnl, unpnl
					coin := CoinAsset{
						CoinName:     sub.Currency.String(),
						Amount:       total,
						FrozenAmount: 0,
					}
					if sub.Currency == goex.USDT {
						coin.Btc = total / btcusdt
						coin.Usdt = total
						coin.Usd = total / usdtusd
						coin.Cny = total / usdtcny
					} else if sub.Currency == goex.BTC {
						coin.Btc = total
						coin.Usdt = total * btcusdt
						coin.Usd = total * btcusd
						coin.Cny = total * btccny
					} else {
						usdt_pair := goex.NewCurrencyPair(sub.Currency, goex.USDT)
						btc_pair := goex.NewCurrencyPair(sub.Currency, goex.BTC)
						usdt_ticker, err := ex.spot.GetTicker(usdt_pair)

						if err != nil {
							btc_ticker, err := ex.spot.GetTicker(btc_pair)
							if err == nil {
								btc := total * btc_ticker.Last
								coin.Btc = btc
								coin.Usdt = btc * btcusdt
								coin.Usd = btc * btcusd
								coin.Cny = btc * btccny
							} else {
								logger.Printf("coin:%s not ticker for value caculate", sub.Currency.String())
							}
						} else {
							usdt := total * usdt_ticker.Last
							coin.Btc = usdt / btcusdt
							coin.Usdt = usdt
							coin.Usd = usdt / usdtusd
							coin.Cny = usdt / usdtcny
						}
					}
					asset, ok := coins[sub.Currency.String()]
					if ok {
						asset.Amount += coin.Amount
						asset.FrozenAmount += coin.FrozenAmount
						asset.Btc += coin.Btc
						asset.Usdt += coin.Usdt
						asset.Usd += coin.Usd
						asset.Cny += coin.Cny
						coins[sub.Currency.String()] = asset
					} else {
						coins[sub.Currency.String()] = coin
					}
				}
			}
		}

		asset := Asset{
			AccountID: ex.accountId,
			Btc:       0,
			Usdt:      0,
			Usd:       0,
			Cny:       0,
			Btc_Usdt:  btcusdt,
			Btc_Usd:   btcusd,
			Btc_Cny:   btccny,
			Usdt_Usd:  usdtusd,
			Usdt_Cny:  usdtcny,
			Usd_Cny:   usdcny,
		}
		coinassets := make([]CoinAsset, 0)
		for _, c := range coins {
			asset.Btc += c.Btc
			asset.Usdt += c.Usdt
			asset.Usd += c.Usd
			asset.Cny += c.Cny
			coinassets = append(coinassets, c)
		}
		orm.AddAsset(asset)
		orm.AddCoinAssets(coinassets)
	}
}

func StartFetchAccount(ctx context.Context, period time.Duration) {
	go NewWorker(ctx, period, UpdateAccounts)
}
