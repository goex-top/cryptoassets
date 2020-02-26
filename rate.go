package main

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/forex"
	"net"
	"net/http"
	"net/url"
	"time"
)

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
	return q.FiftyDayAverage, nil
}

func GetUSDCNY() (float64, error) {
	return GetFinanceFromYahoo("USDCNY=X")
}

func GetBTCUSD() (float64, error) {
	return GetFinanceFromYahoo("BTCUSD=X")
}

func GetBTCUSDCNY() (float64, float64, error) {
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
			usdcny = q.FiftyDayAverage
		} else if q.Symbol == s_btcusd {
			btcusd = q.FiftyDayAverage
		}
	}
	return btcusd, usdcny, nil
}

// https://api.coinmarketcap.com/v2/ticker/825/
//{
//"attention": "WARNING: This API is now deprecated and will be taken offline soon.  Please switch to the new CoinMarketCap API to avoid interruptions in service. (https://pro.coinmarketcap.com/migrate/)",
//"data": {
//"id": 825,
//"name": "Tether",
//"symbol": "USDT",
//"website_slug": "tether",
//"rank": 5,
//"circulating_supply": 4642367414,
//"total_supply": 4776930644,
//"max_supply": null,
//"quotes": {
//"USD": {
//"price": 1.0022976076,
//"volume_24h": 59823709328.9464,
//"market_cap": 4653033752,
//"percent_change_1h": 0.14,
//"percent_change_24h": 0.03,
//"percent_change_7d": 0.11
//}
//},
//"last_updated": 1582689198
//},
//"metadata": {
//"timestamp": 1582688376,
//"warning": "WARNING: This API is now deprecated and will be taken offline soon.  Please switch to the new CoinMarketCap API to avoid interruptions in service. (https://pro.coinmarketcap.com/migrate/)",
//"error": null
//}
//}
