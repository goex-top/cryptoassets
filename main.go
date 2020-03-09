package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/piquette/finance-go"
	"log"
	"os"
	"os/signal"
	"syscall"
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

func usage() {
	fmt.Fprintf(os.Stderr, `version: v1.0.0
Usage:  [-h] [-d=true|false]
Options:
`)
	flag.PrintDefaults()
}

func main() {

	var help bool
	var debug bool
	flag.BoolVar(&debug, "d", false, "enable debug mode")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage

	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	conf, _ = loadConfig()
	orm.DB = initOrm(conf.Debug || debug)

	ctx, cancel = context.WithCancel(context.Background())
	initExchanges(conf)
	initYahooBackend()

	ctx1, cancel1 := context.WithCancel(context.Background())

	UpdateRate()

	StartFetchRate(ctx1)

	StartFetchAccount(ctx, time.Duration(conf.Freq)*time.Second)
	e := echo.New()

	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	if !debug {
		e.File("/", "web/index.html")
		e.Static("static", "web/static")
	}

	route(e)
	if conf.Debug || debug {
		e.Debug = true
	}

	go func() {
		exitSignal := make(chan os.Signal, 1)
		sigs := []os.Signal{os.Interrupt, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM}
		signal.Notify(exitSignal, sigs...)

		s := <-exitSignal
		cancel1()

		fmt.Println("get signal:", s, ", exit!")
		fmt.Println(e.Close())
	}()

	e.Logger.Fatal(e.Start(":9000"))
}
