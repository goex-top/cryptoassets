package main

import (
	"github.com/labstack/echo"
	"github.com/piquette/finance-go"
	"log"
	"os"
)

var (
	conf         Config
	exchanges    []Exchange
	accounts     []Account
	orm          OrmManager
	yahooBackEnd *finance.Backends
	rate         Rate
	logger       *log.Logger
)

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	orm.DB = initOrm()
	conf, _ = loadConfig()

	initExchanges(conf)

	e := echo.New()
	route(e)
	e.Logger.Fatal(e.Start(":1323"))
}
