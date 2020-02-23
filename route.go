package main

import "github.com/labstack/echo"

func route(e *echo.Echo) {
	e.GET("/setting", GetSettings)
	e.POST("/setting", AddSettings)
	e.GET("/asset_history", GetAssetHistory)
	e.GET("/asset", GetCurrentAsset)
	e.GET("/exchange_detail", GetCurrentCoins)
}
