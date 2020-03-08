package main

import "github.com/nntaoli-project/goex"

var (
	List = make(map[string][]string)
)

func init() {
	List["OKEx"] = []string{goex.OKEX_FUTURE, goex.OKEX_SWAP, goex.OKEX_V3}
	List["Huobi"] = []string{goex.HBDM, goex.HUOBI_PRO}
	List["BitMEX"] = []string{goex.BITMEX}
	List["Poloniex"] = []string{goex.POLONIEX}
	List["Bitstamp"] = []string{goex.BITSTAMP}
	List["Binance"] = []string{goex.BINANCE, goex.BINANCE_SWAP}
	List["Bittrex"] = []string{goex.BITTREX}
	List["Bithumb"] = []string{goex.BITHUMB}
	List["Gate.io"] = []string{goex.GATEIO}
	List["ZB"] = []string{goex.ZB}
	List["BigONE"] = []string{goex.BIGONE}
	List["HitBTC"] = []string{goex.HITBTC}
}
