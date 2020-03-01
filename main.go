package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/piquette/finance-go"
	"log"
	"os"
	"time"
)

const tokenSecKey = "crypto_asset_token"

var (
	conf         Config
	exchanges    Exchanges
	accounts     []Account
	orm          OrmManager
	yahooBackEnd *finance.Backends
	rate         Rate
	logger       *log.Logger
	ctx          context.Context
	cancel       func()
)

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	orm.DB = initOrm()
	conf, _ = loadConfig()

	ctx, cancel = context.WithCancel(context.Background())
	initExchanges(conf)
	StartFetchAccount(ctx, time.Duration(conf.Freq)*time.Second)
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	route(e)
	e.Debug = true

	//e.Logger.SetLevel(elog.DEBUG)
	e.Logger.Fatal(e.Start(":9000"))
}
