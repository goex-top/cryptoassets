package main

import (
	"context"
	"errors"
	"github.com/goex-top/market_center"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
	"sync"
	"time"
)

type Exchange struct {
	id        uint
	name      string
	nickname  string
	accountId uint
	spot      goex.API
	future    []goex.FutureRestAPI
}

type Exchanges struct {
	mu        sync.RWMutex
	exchanges []Exchange
}

func initExchanges(config Config) {
	orm.Find(&accounts)

	for _, acc := range accounts {
		addExchange(acc)
	}
}

func addExchange(account Account) {
	secKey := ""
	passKey := ""
	if account.ApiSecretKey != "" {
		secKey = string(AESECBDecrypt([]byte(account.ApiSecretKey), []byte(conf.User.Password)))
	}
	if account.ApiPassphrase != "" {
		passKey = string(AESECBDecrypt([]byte(account.ApiPassphrase), []byte(conf.User.Password)))
	}

	exc, isOk := List[account.ExchangeName]
	if !isOk {
		return
	}
	exchange := Exchange{id: account.ID, name: account.ExchangeName, nickname: account.NickName, accountId: account.ID}
	for _, name := range exc {
		if market_center.IsFutureExchange(name) {
			if passKey != "" {
				exchange.future = append(exchange.future,
					builder.NewAPIBuilder().HttpProxy(conf.Proxy).
						APIKey(account.ApiKey).APISecretkey(secKey).ApiPassphrase(passKey).
						BuildFuture(name))
			} else {
				exchange.future = append(exchange.future,
					builder.NewAPIBuilder().HttpProxy(conf.Proxy).
						APIKey(account.ApiKey).APISecretkey(secKey).
						BuildFuture(name))
			}
		} else {
			if passKey != "" {
				exchange.spot = builder.NewAPIBuilder().HttpProxy(conf.Proxy).
					APIKey(account.ApiKey).APISecretkey(secKey).ApiPassphrase(passKey).Build(name)
			} else {
				exchange.spot = builder.NewAPIBuilder().HttpProxy(conf.Proxy).
					APIKey(account.ApiKey).APISecretkey(secKey).Build(name)
			}
		}
	}
	exchanges.mu.Lock()
	exchanges.exchanges = append(exchanges.exchanges, exchange)
	exchanges.mu.Unlock()
}

func deleteExchange(id uint) {
	exchanges.mu.Lock()
	for k, v := range exchanges.exchanges {
		if v.id == id {
			exchanges.exchanges = append(exchanges.exchanges[:k], exchanges.exchanges[k+1:]...)
			break
		}
	}
	exchanges.mu.Unlock()
}

func verifyAccount(account Account) error {
	secKey := ""
	passKey := ""
	if account.ApiSecretKey != "" {
		secKey = string(AESECBDecrypt([]byte(account.ApiSecretKey), []byte(conf.User.Password)))
	}
	if account.ApiPassphrase != "" {
		passKey = string(AESECBDecrypt([]byte(account.ApiPassphrase), []byte(conf.User.Password)))
	}

	exc, isOk := List[account.ExchangeName]
	if !isOk {
		return errors.New("do not support exchange")
	}

	for _, name := range exc {
		if market_center.IsFutureExchange(name) {
			var ex goex.FutureRestAPI

			if passKey != "" {

				ex = builder.NewAPIBuilder().HttpProxy(conf.Proxy).
					APIKey(account.ApiKey).APISecretkey(secKey).ApiPassphrase(passKey).
					BuildFuture(name)
			} else {
				ex = builder.NewAPIBuilder().HttpProxy(conf.Proxy).
					APIKey(account.ApiKey).APISecretkey(secKey).
					BuildFuture(name)
			}
			_, err := ex.GetFutureUserinfo()
			return err
		} else {
			var ex goex.API
			if passKey != "" {
				ex = builder.NewAPIBuilder().HttpProxy(conf.Proxy).
					APIKey(account.ApiKey).APISecretkey(secKey).ApiPassphrase(passKey).Build(name)
			} else {
				ex = builder.NewAPIBuilder().HttpProxy(conf.Proxy).
					APIKey(account.ApiKey).APISecretkey(secKey).Build(name)
			}
			_, err := ex.GetAccount()
			return err
		}
	}
	return nil
}

func addAccount(account Account) {
	accounts = append(accounts, account)
	addExchange(account)
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
}

func deleteAccount(id uint) {
	for k, v := range accounts {
		if v.ID == id {
			accounts = append(accounts[:k], accounts[k+1:]...)
			break
		}
	}
	deleteExchange(id)
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
}

func UpdateAccounts() {
	usdcny := GetUsdCny()
	btcusd := GetBtcUsd()
	usdtusd := GetUsdtUsd()

	btcusdt := btcusd * usdtusd
	btccny := btcusd * usdcny
	usdtcny := usdtusd * usdcny

	exchanges.mu.RLock()
	exs := exchanges.exchanges
	exchanges.mu.RUnlock()

	for _, ex := range exs {
		coins := make(map[string]CoinAsset)
		if ex.spot != nil {
			acc, err := ex.spot.GetAccount()
			if err != nil {
				continue
			}
			for _, sub := range acc.SubAccounts {
				total := sub.Amount + sub.ForzenAmount
				if total == 0 {
					continue
				}
				coin := CoinAsset{
					CoinName:     sub.Currency.String(),
					Amount:       sub.Amount,
					FrozenAmount: sub.ForzenAmount,
				}
				if sub.Currency == goex.USDT {
					if btcusdt != 0 {
						coin.Btc = total / btcusdt
					}

					coin.Usdt = total
					if usdtusd != 0 {
						coin.Usd = total / usdtusd
					}

					if usdtcny != 0 {
						coin.Cny = total / usdtcny
					}

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
							logger.Printf("coin: [%s] not ticker for value caculate", sub.Currency.String())
						}
					} else {
						usdt := total * usdt_ticker.Last
						if btcusdt != 0 {
							coin.Btc = usdt / btcusdt
						}

						coin.Usdt = usdt

						if usdtusd != 0 {
							coin.Usd = usdt / usdtusd
						}

						if usdtcny != 0 {
							coin.Cny = usdt / usdtcny
						}

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
								logger.Printf("coin: [%s] not ticker for value caculate", sub.Currency.String())
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
		asset = orm.AddAsset(asset)
		for k := range coinassets {
			coinassets[k].AssetID = asset.ID
		}
		orm.AddCoinAssets(coinassets)
	}
}

func StartFetchAccount(ctx context.Context, period time.Duration) {
	go NewWorker(ctx, period, UpdateAccounts)
}
