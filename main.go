package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/piquette/finance-go"
	"log"
	"os"
)

const tokenSecKey = "crypto_asset_token"

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

	// Middleware
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	route(e)
	e.Debug = true

	//e.Logger.SetLevel(elog.DEBUG)
	e.Logger.Fatal(e.Start(":9000"))
}
