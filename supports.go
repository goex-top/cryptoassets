package main

import (
	"github.com/goex-top/market_center"
)

var (
	List = make(map[string][]string)
)

func init() {
	List["OKEx"] = []string{market_center.OKEX, market_center.FUTURE_OKEX, market_center.SWAP_OKEX}
	List["Huobi"] = []string{market_center.HUOBI, market_center.FUTURE_HBDM}
	List["BitMEX"] = []string{market_center.FUTURE_BITMEX}
	List["Bitfinex"] = []string{market_center.BITFINEX}
	List["Poloniex"] = []string{market_center.POLONIEX}
	List["Bitstamp"] = []string{market_center.BITSTAMP}
	List["Binance"] = []string{market_center.BINANCE, market_center.SWAP_BINANCE}
	List["Bittrex"] = []string{market_center.BITTREX}
	List["Bithumb"] = []string{market_center.BITHUMB}
	List["Gate.io"] = []string{market_center.GATEIO}
	List["ZB"] = []string{market_center.ZB}
	List["BigONE"] = []string{market_center.BIGONE}
	List["HitBTC"] = []string{market_center.HITBTC}
}
