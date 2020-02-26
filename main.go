package main

import (
	"github.com/labstack/echo"
)

var (
	conf      Config
	exchanges []Exchange
	accounts  []Account
	orm       OrmManager
)

func main() {
	orm.DB = initOrm()
	conf, _ = loadConfig()
	initExchanges(conf)
	e := echo.New()
	route(e)
	e.Logger.Fatal(e.Start(":1323"))
}
