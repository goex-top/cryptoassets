package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func route(e *echo.Echo) {

	g := e.Group("/api")

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(tokenSecKey),
		TokenLookup: "header:Authorization",
	}))

	g.GET("/support", GetSupport)

	g.GET("/setting", GetSetting)
	g.POST("/setting", AddSetting)
	g.DELETE("/setting/:id", DeleteSetting)

	g.GET("/asset_history", GetAssetHistory)
	g.GET("/asset", GetCurrentAsset)
	g.GET("/exchange_detail", GetCurrentCoins)

	// user
	usr := e.Group("/api/user")
	usr.POST("/login", UserLogin)
	usr.POST("/logout", UserLogout)
	usr.GET("/info", GetUserInfo)

}
