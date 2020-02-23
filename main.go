package main

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var (
	conf      Config
	exchanges []Exchange
	accounts  []Account
	orm       *gorm.DB
)

func main() {
	orm = initOrm()
	conf, _ = loadConfig()
	initExchanges(conf)
	e := echo.New()
	route(e)
	e.Logger.Fatal(e.Start(":1323"))
}
