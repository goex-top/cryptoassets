package main

import (
	"context"
	"encoding/json"
	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/binance"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/forex"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Rate struct {
	mu       sync.RWMutex
	BTC_USD  float64
	USDT_USD float64
	USD_CNY  float64
}

func newHttpClient(proxy string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				if proxy != "" {
					return url.Parse(proxy)
				}
				return nil, nil
			},
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
		},
		Timeout: 10 * time.Second,
	}
}

func initYahooBackend() {
	yahooBackEnd = finance.NewBackends(newHttpClient(conf.Proxy))
}

func GetFinanceFromYahoo(symbol string) (float64, error) {
	q, err := forex.Get(symbol)
	if err != nil {
		return 0, err
	}
	return (q.Bid + q.Ask) / 2, nil
}

func GetUSDCNYFromYahoo() (float64, error) {
	return GetFinanceFromYahoo("USDCNY=X")
}

func GetBTCUSDFromYahoo() (float64, error) {
	return GetFinanceFromYahoo("BTCUSD=X")
}

func GetBTCUSDCNYFromYahoo() (float64, float64, error) {
	s_btcusd := "BTCUSD=X"
	s_usdcny := "USDCNY=X"
	symbols := []string{s_btcusd, s_usdcny}
	iter := forex.List(symbols)

	btcusd := 0.0
	usdcny := 0.0
	for iter.Next() {
		if iter.Err() != nil {
			return 0, 0, iter.Err()
		}

		q := iter.ForexPair()
		if q.Symbol == s_usdcny {
			usdcny = (q.Bid + q.Ask) / 2
		} else if q.Symbol == s_btcusd {
			btcusd = (q.Bid + q.Ask) / 2
		}
	}
	return btcusd, usdcny, nil
}

type CoinMarketCapRsp struct {
	Attention string `json:"attention"`
	Data      struct {
		CirculatingSupply float64 `json:"circulating_supply"`
		ID                int     `json:"id"`
		LastUpdated       int     `json:"last_updated"`
		Name              string  `json:"name"`
		Rank              int     `json:"rank"`
		Symbol            string  `json:"symbol"`
		TotalSupply       float64 `json:"total_supply"`
		WebsiteSlug       string  `json:"website_slug"`
		Quotes            struct {
			USD struct {
				MarketCap        float64 `json:"market_cap"`
				PercentChange1h  float64 `json:"percent_change_1h"`
				PercentChange24h float64 `json:"percent_change_24h"`
				PercentChange7d  float64 `json:"percent_change_7d"`
				Price            float64 `json:"price"`
				Volume24h        float64 `json:"volume_24h"`
			} `json:"USD"`
		} `json:"quotes"`
		//MaxSupply         string  `json:"max_supply"`
	} `json:"data"`
	Metadata struct {
		Error     string `json:"error"`
		Timestamp int    `json:"timestamp"`
		Warning   string `json:"warning"`
	} `json:"metadata"`
}

// deprecated
func GetUSDTUSDFromCoinMarketCap() (float64, error) {
	url := "https://api.coinmarketcap.com/v2/ticker/825"
	client := newHttpClient(conf.Proxy)
	rsp, err := goex.HttpGet5(client, url, nil)

	if err != nil {
		return 0, err
	}

	cm := CoinMarketCapRsp{}
	err = json.Unmarshal(rsp, &cm)
	if err != nil {
		return 0, err
	}
	return cm.Data.Quotes.USD.Price, nil
}

// deprecated
func GetBTCUSDFromCoinMarketCap() (float64, error) {
	url := "https://api.coinmarketcap.com/v2/ticker/1"
	client := newHttpClient(conf.Proxy)
	rsp, err := goex.HttpGet5(client, url, nil)

	if err != nil {
		return 0, err
	}

	cm := CoinMarketCapRsp{}
	err = json.Unmarshal(rsp, &cm)
	if err != nil {
		return 0, err
	}
	return cm.Data.Quotes.USD.Price, nil
}

func GetUSDTUSDFromBinanceUS() (float64, error) {
	var ba = binance.NewWithConfig(
		&goex.APIConfig{
			HttpClient: &http.Client{
				Transport: &http.Transport{
					Proxy: func(req *http.Request) (*url.URL, error) {
						if conf.Proxy != "" {
							return url.Parse("socks5://127.0.0.1:1080")
						}
						return nil, nil
					},
					Dial: (&net.Dialer{
						Timeout: 10 * time.Second,
					}).Dial,
				},
				Timeout: 10 * time.Second,
			},
			Endpoint:     binance.US_API_BASE_URL,
			ApiKey:       "",
			ApiSecretKey: "",
		})
	ticker, err := ba.GetTicker(goex.NewCurrencyPair2("USDT_USD"))

	if err != nil {
		return 0, err
	}
	return ticker.Last, nil
}

func GetBTCUSDFromCBinanceUS() (float64, error) {
	var ba = binance.NewWithConfig(
		&goex.APIConfig{
			HttpClient: &http.Client{
				Transport: &http.Transport{
					Proxy: func(req *http.Request) (*url.URL, error) {
						if conf.Proxy != "" {
							return url.Parse("socks5://127.0.0.1:1080")
						}
						return nil, nil
					},
					Dial: (&net.Dialer{
						Timeout: 10 * time.Second,
					}).Dial,
				},
				Timeout: 10 * time.Second,
			},
			Endpoint:     binance.US_API_BASE_URL,
			ApiKey:       "",
			ApiSecretKey: "",
		})
	ticker, err := ba.GetTicker(goex.BTC_USD)

	if err != nil {
		return 0, err
	}
	return ticker.Last, nil
}

func UpdateRate() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		usdcny, err := GetUSDCNYFromYahoo()
		if err == nil {
			updateUsdCny(usdcny)
		} else {
			logger.Printf("GetUSDCNYFromYahoo err:%v", err)
		}
	}()

	go func() {
		defer wg.Done()
		btcusd, err := GetBTCUSDFromCBinanceUS()
		if err == nil {
			updateBtcUsd(btcusd)
		} else {
			logger.Printf("GetBTCUSDFromCoinMarketCap err:%v", err)
		}
	}()

	go func() {
		defer wg.Done()
		usdtusd, err := GetUSDTUSDFromBinanceUS()
		if err == nil {
			updateUsdtUsd(usdtusd)
		} else {
			logger.Printf("GetUSDTUSDFromCoinMarketCap err:%v", err)
		}
	}()
	wg.Wait()

}

func StartFetchRate(ctx context.Context) {
	go NewWorker(ctx, 2*time.Hour, UpdateRate)
}

func updateBtcUsd(btcusd float64) {
	rate.mu.Lock()
	rate.BTC_USD = btcusd
	rate.mu.Unlock()
}

func updateUsdtUsd(usdtusd float64) {
	rate.mu.Lock()
	rate.USDT_USD = usdtusd
	rate.mu.Unlock()
}

func updateUsdCny(usdcny float64) {
	rate.mu.Lock()
	rate.USD_CNY = usdcny
	rate.mu.Unlock()
}

func GetBtcUsd() float64 {
	rate.mu.RLock()
	defer rate.mu.RUnlock()
	return rate.BTC_USD
}

func GetUsdtUsd() float64 {
	rate.mu.RLock()
	defer rate.mu.RUnlock()
	return rate.USDT_USD
}

func GetUsdCny() float64 {
	rate.mu.RLock()
	defer rate.mu.RUnlock()
	return rate.USD_CNY
}
